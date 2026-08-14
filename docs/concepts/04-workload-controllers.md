# 04. Workload Controller

## 학습 목표

- Deployment, StatefulSet, DaemonSet, Job과 CronJob을 용도에 맞게 선택한다.

| Controller | 용도 | 핵심 성질 |
|---|---|---|
| ReplicaSet | 동일 Pod 개수 유지 | 보통 Deployment가 소유 |
| Deployment | stateless rollout·rollback | 새 ReplicaSet으로 교체 |
| StatefulSet | 순서·고정 이름·개별 PVC | 데이터 복제는 별도 책임 |
| DaemonSet | 대상 Node마다 하나 | CNI, 로그·보안 agent |
| Job | 정해진 완료 횟수 | 실패 재시도와 완료 추적 |
| CronJob | 일정마다 Job 생성 | 동시 실행·history 정책 필요 |

## Deployment 업데이트

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxUnavailable: 0
    maxSurge: 1
```

readiness가 성공해야 새 Pod가 가용한 것으로 계산된다. PDB는 rollout controller를 제한하지 않으므로 Deployment 전략을 별도로 설정한다.

```bash
kubectl rollout status -n shop deployment/gateway
kubectl rollout history -n shop deployment/gateway
kubectl rollout undo -n shop deployment/gateway
kubectl scale -n shop deployment/gateway --replicas=3
```

StatefulSet은 `writer-0` 같은 안정적 이름과 PVC 연결을 제공하지만 single replica DB를 HA로 만들지 않는다. quorum, replication, backup과 restore는 해당 데이터 시스템 설계가 필요하다.

## 체크리스트

- [ ] 일반 HTTP API에는 Deployment를 사용한다.
- [ ] 모든 Node에 필요한 agent에만 DaemonSet을 사용한다.
- [ ] CronJob 작업은 중복 실행돼도 안전하도록 idempotency를 검토한다.

실습: [Lab 02](../../labs/02-deployment-service/README.md), [Lab 05](../../labs/05-job-cronjob/README.md), [Lab 06](../../labs/06-storage-statefulset/README.md)
