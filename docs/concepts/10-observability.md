# 10. Observability와 Alerting

## 학습 목표

- 로그, event, metric, trace를 목적에 맞게 사용한다.
- SLI 기반 경보를 runbook과 연결한다.

| 신호 | 답하는 질문 |
|---|---|
| Event | Kubernetes가 왜 스케줄·pull·mount·kill하지 못했나? |
| Log | 특정 요청과 process에서 무슨 일이 있었나? |
| Metric | 시간에 따라 얼마나 자주·느리게·포화됐나? |
| Trace | MSA 요청이 어느 hop에서 지연·실패했나? |

## 네 가지 기본 신호

- traffic: 요청량
- errors: 실패율
- latency: 특히 tail latency(p95/p99)
- saturation: CPU, memory, queue, connection, disk 등의 한계 접근

Prometheus는 metric을 수집하고 rule을 평가한다. Alertmanager는 중복 제거·grouping·routing을 담당하고 Grafana는 dashboard를 제공한다. dashboard만 있고 alert/runbook이 없으면 사람이 계속 보고 있어야 한다.

## Kubernetes 진단 순서

```bash
kubectl get pod -A
kubectl get event -A --sort-by=.lastTimestamp
kubectl describe pod -n shop POD_NAME
kubectl logs -n shop POD_NAME --all-containers --since=15m
kubectl logs -n shop POD_NAME --previous
kubectl top node
kubectl top pod -A
```

Pod 로그는 node와 Pod 수명에 묶인다. 운영에서는 중앙 수집과 보존 정책이 필요하다. 개인정보·credential을 로그에 쓰지 않고 request ID와 service/version/pod 정보를 구조화한다.

## Alert 원칙

- 사용자 영향 또는 조치가 필요한 조건에 page
- 원인 하나당 수십 개 증상 alert를 만들지 않기
- 지속 시간 `for`로 순간 spike 억제
- annotation에 영향, dashboard, runbook 경로 포함
- alert 자체의 전달 실패도 감시

## 체크리스트

- [ ] metric 이름과 SLI 계산식을 문서화한다.
- [ ] 재시작 횟수만이 아니라 사용자 오류율·지연시간도 경보한다.
- [ ] 모든 page alert에 소유자와 runbook이 있다.

실습: [Lab 10](../../labs/10-observability/README.md)
