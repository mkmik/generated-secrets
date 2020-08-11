// +build !ignore_autogenerated

// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecret) DeepCopyInto(out *GeneratedSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(GeneratedSecretStatus)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecret.
func (in *GeneratedSecret) DeepCopy() *GeneratedSecret {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeneratedSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretList) DeepCopyInto(out *GeneratedSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GeneratedSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretList.
func (in *GeneratedSecretList) DeepCopy() *GeneratedSecretList {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeneratedSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretStatus) DeepCopyInto(out *GeneratedSecretStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretStatus.
func (in *GeneratedSecretStatus) DeepCopy() *GeneratedSecretStatus {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretStatus)
	in.DeepCopyInto(out)
	return out
}
