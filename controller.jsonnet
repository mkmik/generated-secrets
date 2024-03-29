{
  // supported configuration points:
  fields:: {
    appImage: 'ko://mkm.pub/generated-secrets',
    namespace: 'kube-system',
  },

  crd: {
    apiVersion: 'apiextensions.k8s.io/v1',
    kind: 'CustomResourceDefinition',
    metadata: {
      name: 'generatedsecrets.mkm.pub',
      labels: {
        app: 'generated-secrets',
      },
    },
    spec: {
      group: 'mkm.pub',
      names: {
        kind: 'GeneratedSecret',
        listKind: 'GeneratedSecretList',
        plural: 'generatedsecrets',
        singular: 'generatedsecret',
      },
      scope: 'Namespaced',
      versions: [
        {
          name: 'v1alpha1',
          served: true,
          storage: true,
          subresources: {
            status: {},
          },
          schema: {
            openAPIV3Schema: {
              type: 'object',
              properties: {
                spec: {
                  type: 'object',
                  'x-kubernetes-preserve-unknown-fields': true,
                },
                status: {
                  'x-kubernetes-preserve-unknown-fields': true,
                },
              },
            },
          },
        },
      ],
    },
  },
  serviceAccount: {
    apiVersion: 'v1',
    kind: 'ServiceAccount',
    metadata: {
      name: 'generated-secrets-controller',
      namespace: $.fields.namespace,
      labels: {
        app: 'generated-secrets',
      },
    },
  },
  deployment: {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
      name: 'generated-secrets-controller',
      namespace: $.fields.namespace,
      labels: {
        app: 'generated-secrets',
      },
      annotations: {
        'field.knot8.io/appImage': '/spec/template/spec/containers/~{"name":"controller"}/image',
        'knot8.io/original': std.manifestYamlDoc($.fields, quote_keys=false),
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          name: 'generated-secrets-controller',
        },
      },
      template: {
        metadata: {
          labels: {
            name: 'generated-secrets-controller',
          },
        },
        spec: {
          containers: [
            {
              name: 'controller',
              image: $.fields.appImage,
              ports: [
                {
                  containerPort: 8080,
                  name: 'http',
                },
              ],
              resources: {
                requests: {
                  memory: '200Mi',
                  cpu: '0.4',
                },
                limits: {
                  memory: '800Gi',
                  cpu: '2',
                },
              },
              securityContext: {
                readOnlyRootFilesystem: true,
                runAsNonRoot: true,
                runAsUser: 1001,
              },
              stdin: false,
              tty: false,
              volumeMounts: [
                {
                  mountPath: '/tmp',
                  name: 'tmp',
                },
              ],
            },
          ],
          serviceAccount: 'generated-secrets-controller',
          securityContext: {
            fsGroup: 65534,
          },
          volumes: [
            {
              emptyDir: {},
              name: 'tmp',
            },
          ],
        },
      },
    },
  },
  rbac: {
    roles: {
      view: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRole',
        metadata: {
          name: 'generated-secrets-view',
          namespace: $.fields.namespace,
          labels: {
            app: 'generated-secrets',
            'rbac.authorization.k8s.io/aggregate-to-admin': 'true',
            'rbac.authorization.k8s.io/aggregate-to-edit': 'true',
            'rbac.authorization.k8s.io/aggregate-to-view': 'true',
          },
        },
        rules: [
          {
            apiGroups: [
              'mkm.pub',
            ],
            resources: [
              'generatedsecrets',
            ],
            verbs: [
              'get',
              'list',
              'watch',
            ],
          },
        ],
      },
      edit: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRole',
        metadata: {
          name: 'generated-secrets-edit',
          namespace: $.fields.namespace,
          labels: {
            app: 'generated-secrets',
            'rbac.authorization.k8s.io/aggregate-to-admin': 'true',
            'rbac.authorization.k8s.io/aggregate-to-edit': 'true',
          },
        },
        rules: [
          {
            apiGroups: [
              'mkm.pub',
            ],
            resources: [
              'generatedsecrets',
            ],
            verbs: [
              'create',
              'delete',
              'deletecollection',
              'patch',
              'update',
            ],
          },
        ],
      },
      controller: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRole',
        metadata: {
          name: 'generated-secrets:controller',
          namespace: $.fields.namespace,
          labels: {
            app: 'generated-secrets',
          },
        },
        rules: [
          {
            apiGroups: [
              '',
            ],
            resources: [
              'secrets',
            ],
            verbs: [
              'get',
              'list',
              'create',
              'delete',
              'deletecollection',
              'patch',
              'update',
              'watch',
            ],
          },
          {
            apiGroups: [
              '',
            ],
            resources: [
              'events',
            ],
            verbs: [
              'create',
              'patch',
            ],
          },
          {
            apiGroups: [
              'mkm.pub',
            ],
            resources: [
              'generatedsecrets',
            ],
            verbs: [
              'get',
              'list',
              'watch',
            ],
          },
          {
            apiGroups: [
              'mkm.pub',
            ],
            resources: [
              'generatedsecrets/status',
            ],
            verbs: [
              'update',
              'patch',
            ],
          },
        ],
      },
    },
    bindings: {
      controller: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRoleBinding',
        metadata: {
          name: 'generated-secrets:controller',
          namespace: $.fields.namespace,
          labels: {
            app: 'generated-secrets',
          },
        },
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'ClusterRole',
          name: 'generated-secrets:controller',
        },
        subjects: [
          {
            kind: 'ServiceAccount',
            name: 'generated-secrets-controller',
            namespace: $.fields.namespace,
          },
        ],
      },
    },
  },
}
