// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

//go:build gen
// +build gen

package scripts

// prevent go mod tidy from removing the code-generator from go.mod.
import _ "k8s.io/code-generator"
