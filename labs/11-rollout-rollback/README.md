# 11. 롤아웃과 롤백

## 목표

잘못된 이미지를 배포하고 Deployment 상태와 Events로 실패를 판별한 뒤 이전 Revision으로 되돌린다.

## 준비

```bash
make deploy ENV=kind
kubectl rollout status deployment -n shop gateway --timeout=120s
```

## 실행

```bash
kubectl set image deployment/gateway -n shop gateway=registry.invalid/gateway:broken
kubectl rollout status deployment -n shop gateway --timeout=30s || true
kubectl get pods -n shop -l app.kubernetes.io/name=gateway
kubectl describe deployment -n shop gateway
```

## 관찰

```bash
kubectl rollout history deployment -n shop gateway
kubectl get rs -n shop -l app.kubernetes.io/name=gateway
kubectl get events -n shop --sort-by=.lastTimestamp
```

RollingUpdate는 새 Pod가 준비되지 않으면 기존 Ready Pod를 남긴다. `maxUnavailable`, readinessProbe, PDB가 서로 다른 단계에서 가용성을 지킨다.

## 장애와 복구

```bash
kubectl rollout undo deployment -n shop gateway
kubectl rollout status deployment -n shop gateway --timeout=120s
curl --fail --silent http://localhost:18080/api/products || true
```

명령과 판단 기준은 [롤아웃 runbook](../../docs/runbooks/rollout-rollback.md)을 참고한다.

## 정리

복구된 이미지와 두 개의 Ready 복제본을 확인한다.
