// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
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
	// TODO: remove this code since we filter generation on watch
	if gs.Status != nil && gs.Status.ObservedGeneration == gs.Generation {
		log.V(2).Info("Already caught up", "generation", gs.Generation)
		return reconcile.Result{}, nil
	}

	oldSec := &corev1.Secret{}
	if err := r.client.Get(ctx, req.NamespacedName, oldSec); errors.IsNotFound(err) {
		log.Info("secret doesn't already exist, generating")
	}
	log.V(2).Info("Old secret", "secret", oldSec)

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
	sec.Data = oldSec.Data
	if sec.Data == nil {
		sec.Data = map[string][]byte{}
	}
	for k, v := range oldSec.Annotations {
		if strings.HasPrefix(k, "ts.mkmik.github.com/") {
			sec.Annotations[k] = v
		}
	}

	for d, k := range gs.Spec.Data {
		merge(&k, gs.Spec.Default)
		tsAnno := timestampAnnotation(d)

		if sec.Annotations == nil {
			sec.Annotations = map[string]string{}
		}

		if validateSecret(oldSec.Data[d], k) {
			log.Info("secret exists, skipping", "key", d)
			sec.Annotations[tsAnno] = oldSec.Annotations[tsAnno]
			continue
		}
		sec.Data[d], err = generateSecret(k)
		if err != nil {
			return reconcile.Result{}, err
		}
		sec.Annotations[tsAnno] = metav1.Now().UTC().Format(time.RFC3339)
	}
	// client.Patch needs this
	sec.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	log.Info("applying", "secret", sec)
	if err := r.client.Patch(ctx, sec, client.Apply, client.FieldOwner(controllerName)); err != nil {
		return reconcile.Result{}, err
	}
	log.Info("applied", "secret", sec)

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

// merge merges non-zero values from s into d
func merge(d *v1alpha1.GeneratedSecretKey, s *v1alpha1.GeneratedSecretKey) {
	if s == nil {
		return
	}
	if d.Length == 0 {
		d.Length = s.Length
	}
	if d.TTL == "" {
		d.TTL = s.TTL
	}
}

func validateSecret(b []byte, k v1alpha1.GeneratedSecretKey) bool {
	if b == nil {
		return false
	}
	if k.Binary {
		return validateBinary(b, k.Length)
	}
	return validateRandomString(b, k.Length, k.GetAlphabet())
}

func generateSecret(k v1alpha1.GeneratedSecretKey) ([]byte, error) {
	if k.Binary {
		return randomBinary(k.Length)
	}
	return randomString(k.Length, k.GetAlphabet())
}

func validateRandomString(b []byte, n int, alphabet string) bool {
	if len(b) != n {
		log.Info("random string lenght mismatch")
		return false
	}
	if !bytes.ContainsAny(b, alphabet) {
		log.Info("random string alphabet mismatch")
		return false
	}
	return true
}

func validateBinary(b []byte, n int) bool {
	if len(b) != n {
		log.Info("binary secret length mismatch")
	}

	return len(b) == n
}

func randomString(n int, alphabet string) ([]byte, error) {
	b := make([]byte, n)
	bi := big.NewInt(int64(len(alphabet)))
	for i := range b {
		idx, err := rand.Int(rand.Reader, bi)
		if err != nil {
			return nil, err
		}

		b[i] = alphabet[idx.Int64()]
	}
	return b, nil
}

func randomBinary(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}
