# Reliability Model

## 목표와 경계

이 모델은 학습용 MSA의 가용성 사고방식을 정의한다. 단일 laptop의 minikube/kind는 인프라 HA를 제공하지 않으며, 데이터가 메모리에 있으므로 durable order SLO는 없다.

## 예제 SLI/SLO

| 항목 | 정의 | 예제 목표 |
|---|---|---|
| Availability SLI | 유효 요청 중 5xx가 아닌 비율 | 30일 99.9% |
| Latency SLI | gateway 요청 p95 | 500ms 이하 |
| Recovery | 단일 Pod 삭제 후 Ready 복구 | 60초 이내 |

예제 값은 학습 기준이며 business 요구로 검증한 운영 약속이 아니다.

## 장애 도메인

| 계층 | 장애 | 완화 | 남는 위험 |
|---|---|---|---|
| Pod | crash/OOM/bad process | Deployment, probe | 반복되는 코드 오류 |
| Node | reboot/network/disk | replicas, topology spread | capacity 부족 |
| Cluster | control plane/etcd | managed EKS 또는 RKE2 HA | 사이트 전체 장애 |
| Network | DNS/policy/LB | DNS 허용, Gateway HA, runbook | 외부 회선 장애 |
| Dependency | timeout/5xx | deadline, 제한된 retry | retry storm |
| Data | 삭제/corruption | replication, backup, restore | 복구 검증 실패 |

## 정비 가능성

2 replicas와 PDB 1은 한 번에 한 Pod를 Eviction API로 제거할 수 있게 한다. 두 Pod가 같은 Node에 몰리지 않도록 topology spread를 사용한다. minikube 한 Node에서는 `ScheduleAnyway`로 학습 가능성을 유지하므로 실제 node HA는 없다.

HPA, PDB와 Node capacity를 같이 계산한다. PDB가 허용해도 새 Pod를 배치할 CPU/memory가 없으면 drain 후 가용성이 떨어진다.

## RTO/RPO

- local lab RTO: Git에서 cluster/app을 재배포하는 절차 학습
- on-prem control plane: etcd off-site snapshot과 restore rehearsal
- EKS: control plane은 관리되지만 PV, DNS, registry, identity와 외부 DB 복구는 별도
- in-memory order RPO: Pod 재시작 시 전체 손실. 교육용 제한으로 명시

## 검증

- `tests/smoke.sh`: 정상 요청 경로
- `tests/resilience.sh`: Pod 삭제와 node drain
- Lab 12: OOM, DNS, NetworkPolicy, bad image, dependency와 capacity 장애
- runbook: 명령뿐 아니라 복구 후 사용자 경로 확인

개념: [Reliability, HA와 DR](../concepts/11-reliability-ha-dr.md)
