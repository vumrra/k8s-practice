# Runbook — Rollout 실패와 Rollback

## 증상

새 revision Pod가 Ready가 되지 않거나 배포 후 오류율·지연시간이 악화된다.

## 즉시 안전 확인

- 이전 replica가 남았는지와 `maxUnavailable`을 확인한다.
- 데이터 migration이 이전 application과 호환되는지 확인하기 전 rollback하지 않는다.

## 진단

```bash
kubectl rollout status -n NAMESPACE deployment/NAME --timeout=2m
kubectl rollout history -n NAMESPACE deployment/NAME
kubectl get replicaset,pod -n NAMESPACE
kubectl describe deployment -n NAMESPACE NAME
kubectl logs -n NAMESPACE -l app.kubernetes.io/name=NAME --tail=100
```

image pull, probe, ConfigMap/Secret key, resource, admission policy와 API compatibility를 확인한다.

## 복구

```bash
kubectl rollout pause -n NAMESPACE deployment/NAME
kubectl rollout undo -n NAMESPACE deployment/NAME
kubectl rollout status -n NAMESPACE deployment/NAME
```

GitOps이면 Git revision도 함께 되돌려 controller가 실패 상태를 다시 적용하지 않게 한다.

## 검증

새 Pod Ready뿐 아니라 실제 API, dependency, error/latency SLI와 event를 확인한다.

## 예방

immutable image, startup/readiness, canary, backwards-compatible migration, automated smoke와 rollback 기준을 둔다.
