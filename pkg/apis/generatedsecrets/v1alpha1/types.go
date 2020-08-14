// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultLength = 12
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// GeneratedSecret is a blueprint for a generated secret
type GeneratedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GeneratedSecretSpec `json:"spec"`
	// +optional
	Status *GeneratedSecretStatus `json:"status,omitempty"`
}

type GeneratedSecretStatus struct {
	// ObservedGeneration reflects the generation most recently observed by the generatedgroups controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GeneratedSecretList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []GeneratedSecret `json:"items"`
}

type GeneratedSecretSpec struct {
	Data    map[string]GeneratedSecretKey `json:"data,omitempty"`
	Default *GeneratedSecretKey           `json:"default,omitempty"`
	// +optional
	Template *GeneratedSecretTemplate `json:"template,omitempty"`
}

type GeneratedSecretKey struct {
	Length   int    `json:"length,omitempty"`
	Alphabet string `json:"alphabet,omitempty"`
	Binary   bool   `json:"binary,omitempty"`
	TTL      string `json:"ttl,omitempty"`
}

const defaultAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYabcdefghijklmnopqrstuvwxy0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

func (k *GeneratedSecretKey) GetAlphabet() string {
	if k.Alphabet == "" {
		return defaultAlphabet
	}
	return k.Alphabet
}

type GeneratedSecretTemplate struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}
