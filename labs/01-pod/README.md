# Lab 01 — Pod

## 목표

Pod가 Kubernetes의 최소 배포 단위이며, 삭제한 Pod는 스스로 돌아오지 않는다는 점을 확인한다.

## 준비

```bash
minikube start --driver=docker
docker build -f build/Dockerfile --build-arg APP=apps/hello-api --build-arg BINARY=hello-api -t k8s-practice/hello-api:dev .
minikube image load k8s-practice/hello-api:dev
kubectl apply -f deploy/base/hello/namespace.yaml
```

## 실행

```yaml
# /tmp/hello-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-pod
  namespace: practice
  labels:
    app: hello-pod
spec:
  automountServiceAccountToken: false
  securityContext:
    runAsNonRoot: true
    runAsUser: 65532
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: hello-api
      image: k8s-practice/hello-api:dev
      imagePullPolicy: IfNotPresent
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities: {drop: ["ALL"]}
```

```bash
kubectl apply -f /tmp/hello-pod.yaml
kubectl get pod -n practice -o wide
kubectl describe pod -n practice hello-pod
kubectl logs -n practice hello-pod
```

## 장애와 복구

```bash
kubectl delete pod -n practice hello-pod
kubectl get pod -n practice
```

새 Pod가 생기지 않는다. Pod 자체에는 원하는 개수를 유지하는 controller가 없기 때문이다. 파일을 다시 적용해 복구한다.

## 정리

```bash
kubectl delete -f /tmp/hello-pod.yaml --ignore-not-found
```

