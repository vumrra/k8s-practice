# Lab 06 — PV, PVC와 StatefulSet

## 목표

Pod와 데이터 수명주기가 다르며 StatefulSet만으로 데이터 HA가 완성되지 않는다는 점을 확인한다.

## 준비

동적 provisioning을 제공하는 기본 StorageClass가 필요하다. minikube는 기본 storage provisioner를 제공한다.

## 실행

```yaml
# /tmp/stateful.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata: {name: writer, namespace: practice}
spec:
  serviceName: writer
  replicas: 1
  selector: {matchLabels: {app: writer}}
  template:
    metadata: {labels: {app: writer}}
    spec:
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000
        seccompProfile: {type: RuntimeDefault}
      containers:
        - name: writer
          image: busybox:1.36.1
          command: ["sh", "-c", "test -f /data/value || date > /data/value; cat /data/value; sleep 3600"]
          volumeMounts: [{name: data, mountPath: /data}]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities: {drop: ["ALL"]}
  volumeClaimTemplates:
    - metadata: {name: data}
      spec:
        accessModes: ["ReadWriteOnce"]
        resources: {requests: {storage: 100Mi}}
```

```bash
kubectl apply -f /tmp/stateful.yaml
kubectl rollout status -n practice statefulset/writer
kubectl logs -n practice writer-0
kubectl delete pod -n practice writer-0
kubectl wait -n practice --for=condition=ready pod/writer-0 --timeout=90s
kubectl logs -n practice writer-0
```

Pod가 다시 만들어져도 PVC의 값은 유지된다. 그러나 단일 디스크 장애, 백업, 복제와 복구는 StorageClass·CSI·스토리지 시스템의 책임이다.

## 정리

```bash
kubectl delete -f /tmp/stateful.yaml
kubectl delete pvc -n practice data-writer-0
```
