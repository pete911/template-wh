# template-wh
Kubernetes template mutating admission webhook. This is generic webhook, that takes configmap with key value pairs and
uses it to replace placeholders (`$key`) in kubernetes manifests.

## running template-wh

### requirements
 - install template-wh either from [helm chart](https://pete911.github.io/template-wh/) or manually
 - install [cert-manager](https://cert-manager.io/)

### example configuration
create configmap e.g.:
```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: template-wh
  namespace: kube-system
data:
  cluster: minikube
```

create webhook configuration e.g.:
```yaml
---
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: pod-wh
  namespace: kube-system
  annotations:
    cert-manager.io/inject-ca-from: kube-system/template-wh
webhooks:
  - name: template-wh.kube-system.svc
    rules:
      - operations: ["CREATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
    clientConfig:
      service:
        name: template-wh
        namespace: kube-system
        path: /mutate
        port: 443
    admissionReviewVersions: ["v1"]
    sideEffects: None
    timeoutSeconds: 5
```

Every request to create a pod that contains `$cluster` placeholder either in `metadata` or `spec` field, will be
replaced for `minikube`.

This is just example, template-wh can be used on any resource.

## releases

Releases are automated and triggered on [chart version](charts/template-wh/Chart.yaml) update. If there is any change
to Chart.yaml and the change is on the main branch, this will trigger github action which tags the branch and releases
[chart](https://pete911.github.io/template-wh/) and
[docker image](https://hub.docker.com/repository/docker/pete911/template-wh) with this version.
