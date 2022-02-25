#!/bin/bash

set -e

mod=mkm.pub/generated-secrets/

gen_mod=$(go list -m -f '{{.Dir}}' k8s.io/code-generator)
gen_groups_path=$gen_mod/generate-groups.sh

mod_root="$(dirname ${BASH_SOURCE})/.."

bash $gen_groups_path \
	deepcopy $mod/pkg/client $mod/pkg/apis generatedsecrets:v1alpha1 \
	--output-base "${mod_root}" \
	--go-header-file ./scripts/gen-boilerplate.txt

rsync -rv "${mod_root}/${mod}/" "${mod_root}/"
rm -rf "${mod_root}/${mod}/"
