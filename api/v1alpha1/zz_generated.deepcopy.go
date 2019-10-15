// +build !ignore_autogenerated

/*
Copyright 2019 Google LLC.

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

// autogenerated by controller-gen object, do not modify manually

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkCluster) DeepCopyInto(out *FlinkCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkCluster.
func (in *FlinkCluster) DeepCopy() *FlinkCluster {
	if in == nil {
		return nil
	}
	out := new(FlinkCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlinkCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkClusterComponentState) DeepCopyInto(out *FlinkClusterComponentState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkClusterComponentState.
func (in *FlinkClusterComponentState) DeepCopy() *FlinkClusterComponentState {
	if in == nil {
		return nil
	}
	out := new(FlinkClusterComponentState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkClusterComponentsStatus) DeepCopyInto(out *FlinkClusterComponentsStatus) {
	*out = *in
	out.JobManagerDeployment = in.JobManagerDeployment
	out.JobManagerService = in.JobManagerService
	if in.JobManagerIngress != nil {
		in, out := &in.JobManagerIngress, &out.JobManagerIngress
		*out = new(JobManagerIngressStatus)
		(*in).DeepCopyInto(*out)
	}
	out.TaskManagerDeployment = in.TaskManagerDeployment
	if in.Job != nil {
		in, out := &in.Job, &out.Job
		*out = new(JobStatus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkClusterComponentsStatus.
func (in *FlinkClusterComponentsStatus) DeepCopy() *FlinkClusterComponentsStatus {
	if in == nil {
		return nil
	}
	out := new(FlinkClusterComponentsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkClusterList) DeepCopyInto(out *FlinkClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FlinkCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkClusterList.
func (in *FlinkClusterList) DeepCopy() *FlinkClusterList {
	if in == nil {
		return nil
	}
	out := new(FlinkClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlinkClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkClusterSpec) DeepCopyInto(out *FlinkClusterSpec) {
	*out = *in
	in.ImageSpec.DeepCopyInto(&out.ImageSpec)
	in.JobManagerSpec.DeepCopyInto(&out.JobManagerSpec)
	in.TaskManagerSpec.DeepCopyInto(&out.TaskManagerSpec)
	if in.JobSpec != nil {
		in, out := &in.JobSpec, &out.JobSpec
		*out = new(JobSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.FlinkProperties != nil {
		in, out := &in.FlinkProperties, &out.FlinkProperties
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkClusterSpec.
func (in *FlinkClusterSpec) DeepCopy() *FlinkClusterSpec {
	if in == nil {
		return nil
	}
	out := new(FlinkClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkClusterStatus) DeepCopyInto(out *FlinkClusterStatus) {
	*out = *in
	in.Components.DeepCopyInto(&out.Components)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkClusterStatus.
func (in *FlinkClusterStatus) DeepCopy() *FlinkClusterStatus {
	if in == nil {
		return nil
	}
	out := new(FlinkClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageSpec) DeepCopyInto(out *ImageSpec) {
	*out = *in
	if in.PullSecrets != nil {
		in, out := &in.PullSecrets, &out.PullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageSpec.
func (in *ImageSpec) DeepCopy() *ImageSpec {
	if in == nil {
		return nil
	}
	out := new(ImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobManagerIngressSpec) DeepCopyInto(out *JobManagerIngressSpec) {
	*out = *in
	if in.HostFormat != nil {
		in, out := &in.HostFormat, &out.HostFormat
		*out = new(string)
		**out = **in
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.UseTLS != nil {
		in, out := &in.UseTLS, &out.UseTLS
		*out = new(bool)
		**out = **in
	}
	if in.TLSSecretName != nil {
		in, out := &in.TLSSecretName, &out.TLSSecretName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobManagerIngressSpec.
func (in *JobManagerIngressSpec) DeepCopy() *JobManagerIngressSpec {
	if in == nil {
		return nil
	}
	out := new(JobManagerIngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobManagerIngressStatus) DeepCopyInto(out *JobManagerIngressStatus) {
	*out = *in
	if in.URLs != nil {
		in, out := &in.URLs, &out.URLs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobManagerIngressStatus.
func (in *JobManagerIngressStatus) DeepCopy() *JobManagerIngressStatus {
	if in == nil {
		return nil
	}
	out := new(JobManagerIngressStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobManagerPorts) DeepCopyInto(out *JobManagerPorts) {
	*out = *in
	if in.RPC != nil {
		in, out := &in.RPC, &out.RPC
		*out = new(int32)
		**out = **in
	}
	if in.Blob != nil {
		in, out := &in.Blob, &out.Blob
		*out = new(int32)
		**out = **in
	}
	if in.Query != nil {
		in, out := &in.Query, &out.Query
		*out = new(int32)
		**out = **in
	}
	if in.UI != nil {
		in, out := &in.UI, &out.UI
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobManagerPorts.
func (in *JobManagerPorts) DeepCopy() *JobManagerPorts {
	if in == nil {
		return nil
	}
	out := new(JobManagerPorts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobManagerSpec) DeepCopyInto(out *JobManagerSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(JobManagerIngressSpec)
		(*in).DeepCopyInto(*out)
	}
	in.Ports.DeepCopyInto(&out.Ports)
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Mounts != nil {
		in, out := &in.Mounts, &out.Mounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobManagerSpec.
func (in *JobManagerSpec) DeepCopy() *JobManagerSpec {
	if in == nil {
		return nil
	}
	out := new(JobManagerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobSpec) DeepCopyInto(out *JobSpec) {
	*out = *in
	if in.ClassName != nil {
		in, out := &in.ClassName, &out.ClassName
		*out = new(string)
		**out = **in
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Savepoint != nil {
		in, out := &in.Savepoint, &out.Savepoint
		*out = new(string)
		**out = **in
	}
	if in.AllowNonRestoredState != nil {
		in, out := &in.AllowNonRestoredState, &out.AllowNonRestoredState
		*out = new(bool)
		**out = **in
	}
	if in.Parallelism != nil {
		in, out := &in.Parallelism, &out.Parallelism
		*out = new(int32)
		**out = **in
	}
	if in.NoLoggingToStdout != nil {
		in, out := &in.NoLoggingToStdout, &out.NoLoggingToStdout
		*out = new(bool)
		**out = **in
	}
	if in.RestartPolicy != nil {
		in, out := &in.RestartPolicy, &out.RestartPolicy
		*out = new(v1.RestartPolicy)
		**out = **in
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Mounts != nil {
		in, out := &in.Mounts, &out.Mounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobSpec.
func (in *JobSpec) DeepCopy() *JobSpec {
	if in == nil {
		return nil
	}
	out := new(JobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobStatus) DeepCopyInto(out *JobStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobStatus.
func (in *JobStatus) DeepCopy() *JobStatus {
	if in == nil {
		return nil
	}
	out := new(JobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerPorts) DeepCopyInto(out *TaskManagerPorts) {
	*out = *in
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = new(int32)
		**out = **in
	}
	if in.RPC != nil {
		in, out := &in.RPC, &out.RPC
		*out = new(int32)
		**out = **in
	}
	if in.Query != nil {
		in, out := &in.Query, &out.Query
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerPorts.
func (in *TaskManagerPorts) DeepCopy() *TaskManagerPorts {
	if in == nil {
		return nil
	}
	out := new(TaskManagerPorts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerSpec) DeepCopyInto(out *TaskManagerSpec) {
	*out = *in
	in.Ports.DeepCopyInto(&out.Ports)
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Mounts != nil {
		in, out := &in.Mounts, &out.Mounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Sidecars != nil {
		in, out := &in.Sidecars, &out.Sidecars
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerSpec.
func (in *TaskManagerSpec) DeepCopy() *TaskManagerSpec {
	if in == nil {
		return nil
	}
	out := new(TaskManagerSpec)
	in.DeepCopyInto(out)
	return out
}
