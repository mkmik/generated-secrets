
//go:build gen
// +build gen

package scripts

// prevent go mod tidy from removing the code-generator from go.mod.
import _ "k8s.io/code-generator"
