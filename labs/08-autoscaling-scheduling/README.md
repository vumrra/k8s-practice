# 08. HPA와 스케줄링

## 목표

리소스 requests가 스케줄링과 HPA 계산의 기준이 되는 과정을 확인한다.

## 준비

```bash
make deploy ENV=kind
```

[Metrics Server](../../deploy/addons/metrics/README.md)를 설치하고 `kubectl top pods -n shop`이 동작하는지 확인한다.

## 실행

```bash
kubectl autoscale deployment gateway -n shop --cpu=50% --min=2 --max=6
kubectl run load -n shop --image=busybox:1.37 --restart=Never -- \
  sh -c 'while true; do wget -q -O- http://gateway:8080/api/products >/dev/null; done'
kubectl get hpa -n shop -w
```

## 관찰

```bash
kubectl describe hpa -n shop gateway
kubectl get pods -n shop -l app.kubernetes.io/name=gateway -o wide
kubectl get deployment -n shop gateway -o jsonpath='{.spec.template.spec.topologySpreadConstraints}'
```

HPA는 즉시 줄어들지 않는다. 안정화 시간 때문에 부하 제거 후에도 한동안 복제본을 유지한다.

## 장애와 복구

requests를 실제 사용량보다 지나치게 작게 두면 과민 확장하고, 크게 두면 확장이 늦어진다. 조정 전후의 `TARGETS`와 Pending Pod Events를 비교한다.

## 정리

```bash
kubectl delete hpa -n shop gateway
kubectl delete pod -n shop load --ignore-not-found
kubectl scale deployment -n shop gateway --replicas=2
```
