// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

// Generated-secrets is a kubernetes controller that generates a Secret resource containing
// one or more random secrets according to a declarative blueprint defined in a GeneratedSecret CR.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bitnami-labs/flagenv"
	"github.com/mkmik/generated-secrets/pkg/apis/generatedsecrets/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	controllers "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName = "generated-secret"
)

var (
	logger = logf.Log.WithName(controllerName)
	log    = logger.WithName("main")

	scheme = runtime.NewScheme()
)

func init() {
	corev1.AddToScheme(scheme)
	v1alpha1.AddToScheme(scheme)
}

// flags are flags.
type Flags struct {
	Namespace string
}

func (f *Flags) Bind(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.StringVar(&f.Namespace, "namespace", "", "Limit to this namespace.")
}

// mainE is the main function, but which can return an error instead of having to log at every error check.
func mainE(flags Flags) error {
	log.Info("main", "flags", flags)

	restConfig, err := controllers.GetConfig()
	if err != nil {
		return err
	}

	mgr, err := manager.New(restConfig, manager.Options{
		Scheme:    scheme,
		Namespace: flags.Namespace,
	})
	if err != nil {
		return err
	}
	c, err := controller.New(controllerName, mgr, controller.Options{
		Reconciler: &GeneratedSecretReconciler{
			log:    logger.WithName("reconciler"),
			client: mgr.GetClient(),
		},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1alpha1.GeneratedSecret{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return fmt.Errorf("watching: %w", err)
	}

	log.Info("Starting manager")
	return mgr.Start(signals.SetupSignalHandler())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var flags Flags
	flags.Bind(nil)
	klog.InitFlags(nil)
	flagenv.SetFlagsFromEnv("GENERATED_SECRETS", flag.CommandLine)
	flag.Parse()
	logf.SetLogger(klogr.New())

	if err := mainE(flags); err != nil {
		log.Error(err, "main")
		os.Exit(1)
	}
}
