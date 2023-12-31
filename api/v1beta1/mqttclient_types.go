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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientConfig struct {
	ClientCount             int      `json:"clientCount"`
	ClientModelId           string   `json:"clientModelId,omitempty"`
	PublishQoS              string   `json:"publishQoS,omitempty"`
	SubscribeQoS            string   `json:"subscribeQoS,omitempty"`
	MessagePerHourPerClient int      `json:"messagePerHourPerClient,omitempty"`
	PublishTopics           []string `json:"publishTopics,omitempty"`
	SubscribeTopics         []string `json:"subscribeTopics,omitempty"`
	MinCpuRequests          int      `json:"minCpuRequests,omitempty"`
	MinMemoryRequests       int      `json:"minMemoryRequests,omitempty"`
	MaxCpuLimits            int      `json:"maxCpuLimits,omitempty"`
	MaxMemoryLimits         int      `json:"maxMemoryLimits,omitempty"`
}

type SimulationPod struct {
	SimulationPodId                      string   `json:"simulationPodId"`
	ClientCount                          int      `json:"clientCount"`
	ClientModelId                        string   `json:"clientModelId,omitempty"`
	PublishQoS                           string   `json:"publishQoS,omitempty"`
	SubscribeQoS                         string   `json:"subscribeQoS,omitempty"`
	PublishTopics                        []string `json:"publishTopics,omitempty"`
	SubscribeTopics                      []string `json:"subscribeTopics,omitempty"`
	ConnectionLimitAllocatedPerSecond    string   `json:"connectionLimitAllocatedPerSecond,omitempty"`    // TODO float usage is highly discouraged, as support for them varies across languages. Is int ok?
	MessageSendPerHourPerClientRequested string   `json:"messageSendPerHourPerClientRequested,omitempty"` // TODO float usage is highly discouraged, as support for them varies across languages. Is int ok?
	MessageSendPerHourPerClientAllocated string   `json:"messageSendPerHourPerClientAllocated,omitempty"` // TODO float usage is highly discouraged, as support for them varies across languages. Is int ok?
}

// MqttClientSpec defines the desired state of MqttClient
type MqttClientSpec struct {
	Id                       string         `json:"id"`
	ClientImageId            string         `json:"clientImageId"`
	HostName                 string         `json:"hostName"`
	Port                     int            `json:"port"`
	Protocol                 string         `json:"protocol"`
	EnableTls                bool           `json:"enableTls"`
	ConnectionLimitPerSecond int            `json:"connectionLimitPerSecond"`
	SendingLimitPerSecond    int            `json:"sendingLimitPerSecond"`
	ClientConfigs            []ClientConfig `json:"clientConfigs"`
}

// MqttClientStatus defines the observed state of MqttClient
type MqttClientStatus struct {
	SimulationPods []SimulationPod `json:"simulationPods"`
	RunId          string          `json:"runId"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MqttClient is the Schema for the mqttclients API
type MqttClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MqttClientSpec   `json:"spec,omitempty"`
	Status MqttClientStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MqttClientList contains a list of MqttClient
type MqttClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MqttClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MqttClient{}, &MqttClientList{})
}
