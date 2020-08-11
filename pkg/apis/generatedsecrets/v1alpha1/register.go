// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: "mkmik.github.com", Version: "v1alpha1"}
	SchemeBuilder      = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	AddToScheme        = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(&GeneratedSecret{}, &GeneratedSecretList{})
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
