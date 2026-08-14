# Lab 03 — ConfigMap과 Secret

## 목표

이미지와 환경 설정을 분리하고, Secret의 base64가 암호화가 아님을 확인한다.

## 준비

Lab 02의 Hello Deployment가 실행 중이어야 한다.

## 실행

```bash
kubectl apply -k deploy/overlays/minikube
kubectl patch -n practice configmap hello-api --type merge -p '{"data":{"MESSAGE":"changed without rebuilding"}}'
kubectl rollout restart -n practice deployment/hello-api
kubectl rollout status -n practice deployment/hello-api
kubectl port-forward -n practice service/hello-api 8080:80
curl http://127.0.0.1:8080/config
```

환경 변수로 읽는 ConfigMap은 실행 중인 프로세스에 자동 반영되지 않으므로 rollout이 필요하다. volume mount 방식은 파일이 갱신되지만 애플리케이션이 다시 읽어야 한다.

## Secret 확인

```bash
kubectl create secret generic -n practice demo-secret --from-literal=TOKEN=practice-only
kubectl get secret -n practice demo-secret -o jsonpath='{.data.TOKEN}'
kubectl get secret -n practice demo-secret -o jsonpath='{.data.TOKEN}' | base64 --decode
```

base64 값은 누구나 되돌릴 수 있다. 실제 값이 들어간 Secret YAML과 `.env`는 Git에 커밋하지 않는다. 운영 환경에서는 etcd 암호화, 최소 RBAC, 외부 Secret 저장소를 함께 검토한다.

## 정리

```bash
kubectl delete secret -n practice demo-secret
kubectl apply -k deploy/overlays/minikube
```
