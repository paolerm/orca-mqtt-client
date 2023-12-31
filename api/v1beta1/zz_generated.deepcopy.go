//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientConfig) DeepCopyInto(out *ClientConfig) {
	*out = *in
	if in.PublishTopics != nil {
		in, out := &in.PublishTopics, &out.PublishTopics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SubscribeTopics != nil {
		in, out := &in.SubscribeTopics, &out.SubscribeTopics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientConfig.
func (in *ClientConfig) DeepCopy() *ClientConfig {
	if in == nil {
		return nil
	}
	out := new(ClientConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MqttClient) DeepCopyInto(out *MqttClient) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MqttClient.
func (in *MqttClient) DeepCopy() *MqttClient {
	if in == nil {
		return nil
	}
	out := new(MqttClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MqttClient) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MqttClientList) DeepCopyInto(out *MqttClientList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MqttClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MqttClientList.
func (in *MqttClientList) DeepCopy() *MqttClientList {
	if in == nil {
		return nil
	}
	out := new(MqttClientList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MqttClientList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MqttClientSpec) DeepCopyInto(out *MqttClientSpec) {
	*out = *in
	if in.ClientConfigs != nil {
		in, out := &in.ClientConfigs, &out.ClientConfigs
		*out = make([]ClientConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MqttClientSpec.
func (in *MqttClientSpec) DeepCopy() *MqttClientSpec {
	if in == nil {
		return nil
	}
	out := new(MqttClientSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MqttClientStatus) DeepCopyInto(out *MqttClientStatus) {
	*out = *in
	if in.SimulationPods != nil {
		in, out := &in.SimulationPods, &out.SimulationPods
		*out = make([]SimulationPod, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MqttClientStatus.
func (in *MqttClientStatus) DeepCopy() *MqttClientStatus {
	if in == nil {
		return nil
	}
	out := new(MqttClientStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulationPod) DeepCopyInto(out *SimulationPod) {
	*out = *in
	if in.PublishTopics != nil {
		in, out := &in.PublishTopics, &out.PublishTopics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SubscribeTopics != nil {
		in, out := &in.SubscribeTopics, &out.SubscribeTopics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulationPod.
func (in *SimulationPod) DeepCopy() *SimulationPod {
	if in == nil {
		return nil
	}
	out := new(SimulationPod)
	in.DeepCopyInto(out)
	return out
}
