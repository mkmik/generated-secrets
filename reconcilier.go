// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	"github.com/mkmik/generated-secrets/pkg/apis/generatedsecrets/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type GeneratedSecretReconciler struct {
	log    logr.Logger
	client client.Client
}

func (r *GeneratedSecretReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	log := r.log.WithValues("request", req)
	log.Info("Reconciling")

	gs := &v1alpha1.GeneratedSecret{}
	err := r.client.Get(ctx, req.NamespacedName, gs)
	if errors.IsNotFound(err) {
		log.Info("GeneratedSecret delete")
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch GeneratedSecret: %w", err)
	}
	log.V(2).Info("Get", "resource", gs)
	if gs.Status != nil && gs.Status.ObservedGeneration == gs.Generation {
		log.V(2).Info("Already caught up", "generation", gs.Generation)
		return reconcile.Result{}, nil
	}
	name := req.NamespacedName
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: gs.APIVersion,
					Kind:       gs.Kind,
					Name:       gs.Name,
					UID:        gs.UID,
				},
			},
		},
	}
	if gs.Spec.Template != nil {
		sec.Labels = gs.Spec.Template.Labels
		sec.Annotations = gs.Spec.Template.Annotations
	}

	if sec.Data == nil {
		sec.Data = map[string][]byte{}
	}
	for d, k := range gs.Spec.Data {
		merge(&k, gs.Spec.Default)

		sec.Data[d], err = generateSecret(k)
		if err != nil {
			return reconcile.Result{}, err
		}

		if sec.Annotations == nil {
			sec.Annotations = map[string]string{}
		}
		sec.Annotations[timestampAnnotation(d)] = metav1.Now().UTC().Format(time.RFC3339)
	}
	// client.Patch needs this
	sec.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	if err := r.client.Patch(ctx, sec, client.Apply, client.FieldOwner(controllerName)); err != nil {
		return reconcile.Result{}, err
	}
	log.Info("created", "secret", sec)

	if gs.Status == nil {
		gs.Status = &v1alpha1.GeneratedSecretStatus{}
	}
	gs.Status.ObservedGeneration = gs.Generation

	if err := r.client.Status().Update(ctx, gs); err != nil {
		return reconcile.Result{}, fmt.Errorf("cannot update status: %w", err)
	}
	return reconcile.Result{}, nil
}

func timestampAnnotation(d string) string {
	return fmt.Sprintf("ts.mkmik.github.com/Z%sZ", d)
}

func merge(k *v1alpha1.GeneratedSecretKey, d *v1alpha1.GeneratedSecretKey) {
	if d == nil {
		return
	}
	if k.Length == 0 {
		k.Length = d.Length
	}
	if k.TTL == "" {
		k.TTL = d.TTL
	}
}

func generateSecret(k v1alpha1.GeneratedSecretKey) ([]byte, error) {
	s, err := randomString(k.Length)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func randomString(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
