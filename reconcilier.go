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

const (
	annotationPrefix = "ts.mkmik.github.com/"
)

type GeneratedSecretReconciler struct {
	log    logr.Logger
	client client.Client
}

func lset(m *map[string]string, key, value string) {
	if *m == nil {
		*m = map[string]string{}
	}
	(*m)[key] = value
}

func (r *GeneratedSecretReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
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
	for k, v := range oldSec.Annotations {
		if strings.HasPrefix(k, annotationPrefix) {
			lset(&sec.Annotations, k, v)
		}
	}

	for d, k := range gs.Spec.Data {
		setDefaults(&k, gs.Spec.Default)
		if err := updateSecret(d, sec, k); err != nil {
			return reconcile.Result{}, err
		}
	}
	// client.Patch needs this
	sec.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	log.V(2).Info("applying", "secret", sec)
	if err := r.client.Patch(ctx, sec, client.Apply, client.FieldOwner(controllerName)); err != nil {
		return reconcile.Result{}, err
	}
	log.V(2).Info("applied", "secret", sec)

	if gs.Status == nil {
		gs.Status = &v1alpha1.GeneratedSecretStatus{}
	}
	gs.Status.ObservedGeneration = gs.Generation

	log.Info("done")
	if err := r.client.Status().Update(ctx, gs); err != nil {
		return reconcile.Result{}, fmt.Errorf("cannot update status: %w", err)
	}
	return reconcile.Result{}, nil
}

func setDefaults(cfg *v1alpha1.GeneratedSecretKey, def *v1alpha1.GeneratedSecretKey) {
	merge(cfg, def)
	if cfg.Length == 0 {
		cfg.Length = v1alpha1.DefaultLength
	}
}

func updateSecret(d string, sec *corev1.Secret, cfg v1alpha1.GeneratedSecretKey) error {
	tsAnno := metadataAnnotationName(d)

	exists, err := validateSecret(sec.Data[d], sec.Annotations[tsAnno], cfg)
	if err != nil {
		return err
	}
	if exists {
		log.Info("secret exists, skipping", "key", d)
		lset(&sec.Annotations, tsAnno, sec.Annotations[tsAnno])
		return nil
	}

	if sec.Data == nil {
		sec.Data = map[string][]byte{}
	}
	sec.Data[d], err = generateSecret(cfg)
	if err != nil {
		return err
	}

	lset(&sec.Annotations, tsAnno, metav1.Now().UTC().Format(time.RFC3339))
	return nil
}

func metadataAnnotationName(d string) string {
	return fmt.Sprintf("%sZ%sZ", annotationPrefix, d)
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

func validateSecret(b []byte, oldTs string, k v1alpha1.GeneratedSecretKey) (bool, error) {
	if b == nil {
		return false, nil
	}

	if k.TTL != "" {
		if oldTs == "" {
			return false, nil
		}
		ot, err := time.Parse(time.RFC3339, oldTs)
		if err != nil {
			return false, err
		}

		ttl, err := time.ParseDuration(k.TTL)
		if err != nil {
			return false, err
		}

		if time.Since(ot) > ttl {
			return false, nil
		}
	}

	if k.Binary {
		return validateBinary(b, k.Length), nil
	}
	return validateRandomString(b, k.Length, k.GetAlphabet()), nil
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
