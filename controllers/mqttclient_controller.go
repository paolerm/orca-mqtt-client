/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"

	orcav1beta1 "github.com/paolerm/orca-mqtt-client/api/v1beta1"
)

// MqttClientReconciler reconciles a MqttClient object
type MqttClientReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const mqttClientFinalizer = "paermini.com/mqtt-client-finalizer"
const maxMqttClientPerPods = 100

//+kubebuilder:rbac:groups=orca.paermini.com,resources=mqttclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=orca.paermini.com,resources=mqttclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=orca.paermini.com,resources=mqttclients/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MqttClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	mqttClient := &orcav1beta1.MqttClient{}
	err := r.Get(ctx, req.NamespacedName, mqttClient)
	if err != nil {
		return ctrl.Result{}, err
	}

	isCrDeleted := mqttClient.GetDeletionTimestamp() != nil
	if isCrDeleted {
		if controllerutil.ContainsFinalizer(mqttClient, mqttClientFinalizer) {
			// Run finalization logic. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeMqttClient(ctx, req, mqttClient); err != nil {
				return ctrl.Result{}, err
			}

			// Remove mqttClientFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(mqttClient, mqttClientFinalizer)
			err := r.Update(ctx, mqttClient)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// setup status section
	mqttClient.Status = setupStatus(mqttClient.Spec)
	err = r.Status().Update(ctx, mqttClient)
	if err != nil {
		logger.Error(err, "Failed to update CR!")
		return ctrl.Result{}, err
	}

	statefulSetName := mqttClient.Spec.Id
	existingStatefulSet := &appsv1.StatefulSet{}
	ssNamespacedName := types.NamespacedName{
		Name:      statefulSetName,
		Namespace: req.NamespacedName.Namespace,
	}
	logger.Info("Getting statefulSet under namespace " + req.NamespacedName.Namespace + " and name " + statefulSetName + "...")
	err = r.Get(ctx, ssNamespacedName, existingStatefulSet)
	if err != nil {
		logger.Info("Creating statefulSet under namespace " + req.NamespacedName.Namespace + " and name " + statefulSetName + "...")
		// TODO: define statefulSet format
		statefulSet := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      statefulSetName,
				Namespace: req.NamespacedName.Namespace,
				Labels: map[string]string{
					"simulation": statefulSetName,
				},
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"simulation": statefulSetName,
					},
				},
				ServiceName: statefulSetName,
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"simulation": statefulSetName,
						},
						Annotations: map[string]string{
							"prometheus.io/scrape": "true",
						},
					},
					Spec: apiv1.PodSpec{
						ImagePullSecrets: []apiv1.LocalObjectReference{
							{Name: "docker-secret"},
							{Name: "cdpx-acr-secret"},
						},
						Containers: []apiv1.Container{
							{
								Name:            statefulSetName,
								Image:           "paerminiacrplayground.azurecr.io/hello-world:latest", // TODO: custom image
								ImagePullPolicy: "Always",
								Env:             []apiv1.EnvVar{}, // TODO
								// Resources: {},
								// VolumeMounts: {},
								// LivenessProbe: {}.
								// ReadinessProbe: {},
								// TerminationMessagePath: {},
								// TerminationMessagePolicy: File,
							},
						},
						RestartPolicy: "Always",
						// TerminationGracePeriodSeconds: 30,
					},
				},
			},
		}

		err = r.Create(ctx, statefulSet)
		if err != nil {
			logger.Error(err, "Failed to create statefulSet under namespace "+statefulSet.ObjectMeta.Namespace+" with name "+statefulSet.ObjectMeta.Name)
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Updating statefulSet under namespace " + existingStatefulSet.ObjectMeta.Namespace + " and name " + existingStatefulSet.ObjectMeta.Name + "...")
		// TODO: update obj
		err = r.Update(ctx, existingStatefulSet)
		if err != nil {
			logger.Error(err, "Failed to update statefulSet under namespace "+existingStatefulSet.ObjectMeta.Namespace+" and name "+existingStatefulSet.ObjectMeta.Name+"...")
			return ctrl.Result{}, err
		}
	}

	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(mqttClient, mqttClientFinalizer) {
		controllerutil.AddFinalizer(mqttClient, mqttClientFinalizer)
		err = r.Update(ctx, mqttClient)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MqttClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&orcav1beta1.MqttClient{}).
		Complete(r)
}

func (r *MqttClientReconciler) finalizeMqttClient(ctx context.Context, req ctrl.Request, mqttClient *orcav1beta1.MqttClient) error {
	logger := log.FromContext(ctx)
	mqttClientNamePrefix := mqttClient.Spec.Id

	statefulSetList := &appsv1.StatefulSetList{}
	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
		client.MatchingLabels{"simulation": mqttClientNamePrefix},
	}

	err := r.List(ctx, statefulSetList, opts...)
	if err != nil {
		// TODO: ignore not found error
		logger.Error(err, "Failed to get statefulset list!")
		return err
	}

	logger.Info("CR with name " + req.NamespacedName.Name + " under namespace " + req.NamespacedName.Namespace + " marked as deleted. Deleting statefulSet...")
	for i := 0; i < len(statefulSetList.Items); i++ {
		err = r.Delete(ctx, &statefulSetList.Items[i])
		if err != nil {
			logger.Error(err, "Failed to delete statefulset!")
			return err
		}
	}

	logger.Info("Successfully finalized")
	return nil
}

func setupStatus(spec orcav1beta1.MqttClientSpec) orcav1beta1.MqttClientStatus {
	// TODO calculate ConnectionLimitAllocatedPerSecond, MessageSendPerHourPerClientRequested and MessageSendPerHourPerClientAllocated
	simulationPods := []orcav1beta1.SimulationPod{}

	simulationPodId := 0
	for i := 0; i < len(spec.ClientConfigs); i++ {
		clientConfig := spec.ClientConfigs[i]

		totalRemaining := clientConfig.ClientCount
		for totalRemaining > 0 {
			simulationPod := orcav1beta1.SimulationPod{
				SimulationPodId: strconv.Itoa(simulationPodId),
				ClientModelId:   clientConfig.ClientModelId,
				PublishQoS:      clientConfig.PublishQoS,
				SubscribeQoS:    clientConfig.SubscribeQoS,
				PublishTopics:   clientConfig.PublishTopics,
				SubscribeTopics: clientConfig.SubscribeTopics,
			}

			if totalRemaining <= maxMqttClientPerPods {
				simulationPod.ClientCount = totalRemaining
				totalRemaining = 0
			} else {
				simulationPod.ClientCount = maxMqttClientPerPods
				totalRemaining = totalRemaining - maxMqttClientPerPods
			}

			simulationPods = append(simulationPods, simulationPod)
			simulationPodId = simulationPodId + 1
		}
	}

	result := orcav1beta1.MqttClientStatus{
		SimulationPods: simulationPods,
	}

	return result
}
