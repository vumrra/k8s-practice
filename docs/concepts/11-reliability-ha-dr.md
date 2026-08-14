# 11. Reliability, HA와 DR

## 학습 목표

- HA, backup과 disaster recovery를 구분한다.
- cluster·node·workload·data·dependency 계층별 장애를 설계한다.
- SLI/SLO, RTO/RPO와 PDB의 한계를 설명한다.

## 가용성 계층

```text
Cluster HA    다중 control plane, etcd quorum, API load balancer
Node HA       여러 worker와 node/zone 장애 도메인
Workload HA   replicas, readiness, topology spread, PDB, graceful shutdown
Data HA       replication, consistency, quorum, snapshot, backup, restore
Dependency HA timeout, 제한된 retry, idempotency, circuit breaking
DR            클러스터/사이트 소실 후 Git과 backup으로 재구축
```

한 계층만 복제해도 전체 요청 경로의 다른 단일 장애 지점은 남는다. Service Mesh와 StatefulSet은 도구이며 HA 자체가 아니다.

## SLI, SLO와 Error Budget

- SLI: 실제 측정값. 예: 성공 응답 비율, p95 latency
- SLO: 목표. 예: 30일 성공률 99.9%
- error budget: 목표가 허용하는 실패량

예제 MSA는 학습용으로 `5xx가 아닌 응답 / 전체 유효 요청`을 availability SLI로 본다. 운영에서는 client 관점, 제외할 상태, 측정 창과 데이터 누락 처리까지 명시한다.

## RTO와 RPO

- RTO: 장애 후 서비스를 복구하기까지 허용 시간
- RPO: 복구 시점에서 허용 가능한 데이터 손실 시간 범위

replication은 backup을 대체하지 않는다. 잘못된 삭제와 corruption도 복제될 수 있다. backup은 다른 장애 도메인에 보관하고 restore rehearsal로 실제 복구 시간과 무결성을 확인한다.

## Workload HA 체크

- replicas 2 이상과 충분한 Node capacity
- node와 zone topology spread
- readiness와 endpoint 제거
- RollingUpdate `maxUnavailable`/`maxSurge`
- SIGTERM 처리와 종료 유예
- PDB와 drain 가능성
- dependency timeout과 non-idempotent 요청 retry 금지

PDB는 Eviction API를 이용한 자발적 중단만 제한한다. 직접 Pod/Deployment 삭제, rollout, Node 고장과 network partition을 막지 못한다. PDB가 너무 엄격하면 maintenance drain이 영구히 막힐 수 있다.

## 장애 훈련

```bash
kubectl delete pod -n shop -l app.kubernetes.io/name=gateway
kubectl drain NODE --ignore-daemonsets --delete-emptydir-data --timeout=120s
kubectl get pdb -n shop
kubectl rollout status -n shop deployment/gateway
```

자발적 drain과 갑작스러운 Node 중단을 별개로 시험한다. DNS, NetworkPolicy, OOM, bad image, dependency timeout과 capacity 부족도 포함한다.

## 체크리스트

- [ ] 모든 중요 서비스의 장애 도메인과 단일 장애 지점을 적었다.
- [ ] SLO, RTO, RPO가 business 요구와 연결된다.
- [ ] backup restore를 정기적으로 실제 수행한다.
- [ ] PDB가 보호하지 않는 장애를 안다.

실습: [Lab 12](../../labs/12-failure-recovery/README.md), [신뢰성 모델](../architecture/reliability-model.md) · [공식 Pod 중단 문서](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)
