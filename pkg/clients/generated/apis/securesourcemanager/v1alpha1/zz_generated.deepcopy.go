//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	k8sv1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/k8s/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceHostConfigStatus) DeepCopyInto(out *InstanceHostConfigStatus) {
	*out = *in
	if in.Api != nil {
		in, out := &in.Api, &out.Api
		*out = new(string)
		**out = **in
	}
	if in.GitHttp != nil {
		in, out := &in.GitHttp, &out.GitHttp
		*out = new(string)
		**out = **in
	}
	if in.GitSsh != nil {
		in, out := &in.GitSsh, &out.GitSsh
		*out = new(string)
		**out = **in
	}
	if in.Html != nil {
		in, out := &in.Html, &out.Html
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceHostConfigStatus.
func (in *InstanceHostConfigStatus) DeepCopy() *InstanceHostConfigStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceHostConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceObservedStateStatus) DeepCopyInto(out *InstanceObservedStateStatus) {
	*out = *in
	if in.HostConfig != nil {
		in, out := &in.HostConfig, &out.HostConfig
		*out = new(InstanceHostConfigStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
	if in.StateNote != nil {
		in, out := &in.StateNote, &out.StateNote
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceObservedStateStatus.
func (in *InstanceObservedStateStatus) DeepCopy() *InstanceObservedStateStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceObservedStateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecureSourceManagerInstance) DeepCopyInto(out *SecureSourceManagerInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecureSourceManagerInstance.
func (in *SecureSourceManagerInstance) DeepCopy() *SecureSourceManagerInstance {
	if in == nil {
		return nil
	}
	out := new(SecureSourceManagerInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecureSourceManagerInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecureSourceManagerInstanceList) DeepCopyInto(out *SecureSourceManagerInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecureSourceManagerInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecureSourceManagerInstanceList.
func (in *SecureSourceManagerInstanceList) DeepCopy() *SecureSourceManagerInstanceList {
	if in == nil {
		return nil
	}
	out := new(SecureSourceManagerInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecureSourceManagerInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecureSourceManagerInstanceSpec) DeepCopyInto(out *SecureSourceManagerInstanceSpec) {
	*out = *in
	if in.KmsKey != nil {
		in, out := &in.KmsKey, &out.KmsKey
		*out = new(string)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ProjectRef = in.ProjectRef
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecureSourceManagerInstanceSpec.
func (in *SecureSourceManagerInstanceSpec) DeepCopy() *SecureSourceManagerInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(SecureSourceManagerInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecureSourceManagerInstanceStatus) DeepCopyInto(out *SecureSourceManagerInstanceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]k8sv1alpha1.Condition, len(*in))
		copy(*out, *in)
	}
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.ObservedState != nil {
		in, out := &in.ObservedState, &out.ObservedState
		*out = new(InstanceObservedStateStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecureSourceManagerInstanceStatus.
func (in *SecureSourceManagerInstanceStatus) DeepCopy() *SecureSourceManagerInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(SecureSourceManagerInstanceStatus)
	in.DeepCopyInto(out)
	return out
}
