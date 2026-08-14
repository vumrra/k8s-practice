# 07. 리소스, 스케줄링과 Autoscaling

## 학습 목표

- requests·limits, QoS와 eviction 관계를 설명한다.
- Pod 배치와 Pod/Node 확장을 구분한다.

## 리소스

- request: scheduler가 Node 수용 가능 여부를 계산하고 HPA CPU 사용률의 기준이 된다.
- CPU limit: 초과 사용을 throttling한다.
- memory limit: 초과 시 OOM kill될 수 있다.

| QoS | 조건 | 자원 압박 시 일반적 우선순위 |
|---|---|---|
| Guaranteed | 모든 container CPU·memory request=limit | 가장 늦게 eviction |
| Burstable | 일부 request 존재 | 중간 |
| BestEffort | request/limit 없음 | 먼저 eviction |

limit를 무조건 작게 잡으면 안정성이 아니라 throttling과 OOM을 만든다. 관측값과 부하 시험으로 조정한다.

## 배치

- nodeSelector/node affinity: 특정 Node 특성 선택
- pod affinity/anti-affinity: 다른 Pod와 함께 또는 분리
- topology spread: node·zone 같은 장애 도메인에 균등 배치
- taint/toleration: Node가 받아들일 workload 제한
- PriorityClass/preemption: 부족할 때 낮은 우선순위 Pod 축출 가능

hard constraint가 많으면 가용 자원이 있어도 Pending이 된다. HA 목적에는 hostname/zone topology spread와 충분한 Node가 함께 필요하다.

topology spread는 Pod가 스케줄되는 순간의 배치를 결정한다. 장애 Node가 돌아와도 이미 실행 중인 Pod를 자동으로 재배치하지 않으므로, 복구 후 쏠림을 확인하고 필요하면 PDB를 지키는 controlled rollout이나 검증된 Descheduler 정책을 사용한다.

## 확장 계층

```text
HPA -> Pod replicas 변경
VPA -> requests 추천/변경, 재시작 영향 검토
Node autoscaler/Karpenter -> 부족한 Node capacity 추가
```

HPA가 Pod를 늘려도 Node에 여유가 없으면 Pending이다. scale-out 시간 동안 버틸 headroom도 필요하다.

ResourceQuota는 namespace 총량, LimitRange는 object 기본값·범위를 제한한다. 여러 팀이 공유하는 클러스터에서 noisy neighbor와 무제한 사용을 줄인다.

```bash
kubectl top pod -n shop
kubectl describe pod -n shop POD_NAME
kubectl get event -A --sort-by=.lastTimestamp
kubectl get hpa -A
```

## 체크리스트

- [ ] request 없이 HPA CPU 비율을 신뢰하지 않는다.
- [ ] PDB와 HPA replicas, drain 가능성을 함께 계산한다.
- [ ] Pending 원인을 이미지 문제가 아닌 scheduler event에서 먼저 확인한다.

실습: [Lab 04](../../labs/04-probe-resource/README.md), [Lab 08](../../labs/08-autoscaling-scheduling/README.md)
