# Lab 04 — Probe와 리소스

## 목표

startup·liveness·readiness probe의 역할과 requests·limits가 스케줄링과 런타임에 미치는 영향을 구분한다.

## 준비

Hello 이미지를 minikube에 로드하고 `deploy/overlays/minikube`를 적용한다. `kubectl top`까지 실행하려면 Lab 08의 Metrics Server가 필요하다.

## 관찰

```bash
kubectl apply -k deploy/overlays/minikube
kubectl describe -n practice deployment/hello-api
kubectl top pod -n practice
kubectl get pod -n practice -o custom-columns=NAME:.metadata.name,QOS:.status.qosClass
```

Metrics Server가 없다면 `kubectl top`만 실패한다. requests는 스케줄링과 HPA 계산 기준이고, CPU limit은 throttling, memory limit 초과는 OOMKilled로 이어질 수 있다.

## readiness 장애 유도

```bash
kubectl patch -n practice deployment hello-api --type strategic -p '{"spec":{"template":{"spec":{"containers":[{"name":"hello-api","readinessProbe":{"httpGet":{"path":"/readyz","port":9999}}}]}}}}'
kubectl get pod -n practice -w
kubectl get endpointslice -n practice -l kubernetes.io/service-name=hello-api
```

프로세스는 실행 중이어도 Ready가 아니므로 Service endpoint에서 제외된다. liveness를 같은 용도로 사용하면 불필요한 재시작이 반복된다.

## 복구와 정리

```bash
kubectl apply -k deploy/overlays/minikube
kubectl rollout status -n practice deployment/hello-api
```
