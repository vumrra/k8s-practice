# Runbook — OOMKilled, CPU Throttling과 Eviction

## 증상

container 종료 reason이 OOMKilled이거나 Node pressure로 Pod가 Evicted되고 latency가 증가한다.

## 즉시 안전 확인

- replicas와 트래픽 여유를 확인하고 같은 limit으로 무작정 재시작하지 않는다.
- 데이터 process라면 강제 종료 전 consistency와 volume 상태를 확인한다.

## 진단

```bash
kubectl get pod -A -o wide
kubectl describe pod -n NAMESPACE POD
kubectl top pod -n NAMESPACE --containers
kubectl top node
kubectl describe node NODE
kubectl get event -A --field-selector reason=Evicted
```

OOMKilled은 working set, limit, leak와 burst를 구분한다. CPU throttling은 사용률만 낮게 보일 수 있으므로 throttled seconds와 latency를 함께 본다. Eviction은 memory/disk/PID pressure와 QoS를 확인한다.

## 복구

- 확인된 정상 사용량에 맞춰 request/limit 조정
- leak 또는 unbounded queue 수정 후 새 image 배포
- Node disk pressure의 불필요 image/log를 안전하게 정리하고 보존 정책 수정
- capacity가 부족하면 Node 추가 또는 workload 축소

## 검증

재시작·eviction이 멈추고 latency/error SLI가 회복되는지 확인한다. HPA가 request 기준으로 정상 계산되는지도 본다.

## 예방

부하 시험, memory profile, request/limit dashboard, Node pressure alert, log rotation과 capacity headroom을 둔다.
