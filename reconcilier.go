// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/mkmik/generated-secrets/pkg/apis/generatedsecrets/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	log.Info("Generating secret... TODO")

	if gs.Status == nil {
		gs.Status = &v1alpha1.GeneratedSecretStatus{}
	}
	gs.Status.ObservedGeneration = gs.Generation

	if err := r.client.Status().Update(ctx, gs); err != nil {
		return reconcile.Result{}, fmt.Errorf("cannot update status: %w", err)
	}
	return reconcile.Result{}, nil
}
