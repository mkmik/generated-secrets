[![Go](https://github.com/mkmik/generated-secrets/actions/workflows/go.yml/badge.svg)](https://github.com/mkmik/generated-secrets/actions/workflows/go.yml)
[![](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)](https://pkg.go.dev/mkm.pub/generated-secrets?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/mkm.pub/generated-secrets)](https://goreportcard.com/report/mkm.pub/generated-secrets)
![](https://mkm.pub/generated-secrets/workflows/CI/badge.svg)

# Generated Secrets

This project implements a small controller for Kubernetes that generates Secrets from a declarative blue print custom resource.

This project is experimental.

## Documentation

TODO

### Example

```console
$ cat test-gen.yaml
apiVersion: mkm.pub/v1alpha1
kind: GeneratedSecret
metadata:
  name: test
spec:
  data:
    foo:
      length: 12

  template:
    metadata:
      labels:
        foo: bar

$ kubectl apply -f test-gen.yaml
generatedsecret.mkm.pub/test configured

$ kubectl get secret test -o yaml
apiVersion: v1
data:
  foo: MmEwOTQxNTZjYmE0ZDA5ZTM4Y2UwYzE0
kind: Secret
metadata:
  creationTimestamp: "2020-08-11T18:01:10Z"
  labels:
    foo: bar
  name: test
  namespace: default
  ownerReferences:
  - apiVersion: mkm.pub/v1alpha1
    kind: GeneratedSecret
    name: test
    uid: d98eab26-1a46-442f-b381-21276be65d64
  resourceVersion: "9029122"
  selfLink: /api/v1/namespaces/default/secrets/test
  uid: 7a1383db-0e4a-47f8-a446-a94c7bfb7a8e
type: Opaque

$ kubectl delete generatedsecret test
generatedsecret.mkm.pub "test" deleted

$ kubectl get secret test
Error from server (NotFound): secrets "test" not found
```
## Contributing

The go-yaml-edit project team welcomes contributions from the community. Before you start working with generated-secrets, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

generated-secrets is available under the [BSD-2 license](LICENSE).
