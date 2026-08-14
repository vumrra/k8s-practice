# Lab 02 — Deployment와 Service

## 목표

Deployment의 자기 복구와 Service의 고정 DNS·가상 IP를 확인한다.

## 준비

Lab 01의 minikube 클러스터와 로컬 이미지를 사용한다.

## 실행

```bash
kubectl apply -k deploy/overlays/minikube
kubectl rollout status -n practice deployment/hello-api
kubectl get deployment,replicaset,pod,service,endpointslice -n practice
kubectl port-forward -n practice service/hello-api 8080:80
curl http://127.0.0.1:8080/
```

다른 터미널에서 Pod 하나를 삭제한다.

```bash
kubectl delete pod -n practice -l app.kubernetes.io/name=hello-api --wait=false
kubectl get pod -n practice -w
```

ReplicaSet이 새 Pod를 생성하고 Service가 준비된 Pod만 endpoint로 사용한다.

## 확장과 복구

```bash
kubectl scale -n practice deployment/hello-api --replicas=3
kubectl get pod -n practice
kubectl apply -k deploy/overlays/minikube
```

마지막 명령은 Git에 선언된 replicas 2로 되돌린다.

## 정리

```bash
kubectl delete -k deploy/overlays/minikube
```
