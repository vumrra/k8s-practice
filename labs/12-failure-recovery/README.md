# 12. 장애 주입과 복구

## 목표

Pod, 노드, 메모리, DNS, 정책, 의존 서비스, 이미지, PDB, 용량 장애를 같은 진단 순서로 다룬다.

## 준비

```bash
make deploy ENV=kind
kubectl get pods -n shop -o wide
```

운영 환경이 아닌 실습 클러스터에서만 수행한다.

## 실행

한 번에 하나만 주입하고 복구한 뒤 다음 항목으로 넘어간다.

| 상황 | 주입 | 관찰 | 복구 |
|---|---|---|---|
| Pod 삭제 | `kubectl delete pod -n shop -l app.kubernetes.io/name=gateway` | ReplicaSet 재생성 | rollout 완료 대기 |
| 노드 drain | `kubectl drain k8s-practice-worker --ignore-daemonsets --delete-emptydir-data` | Eviction, PDB, 재스케줄 | `kubectl uncordon k8s-practice-worker` |
| OOM | memory limit를 낮춘 뒤 부하 | `OOMKilled`, restart | 원래 매니페스트 재적용 |
| DNS | 잘못된 서비스 이름으로 요청 | `nslookup`, timeout | 서비스 이름 복구 |
| NetworkPolicy | 허용 정책의 selector 변경 | 연결 timeout | Git의 정책 재적용 |
| 의존성 timeout | orders의 inventory 주소를 무응답 주소로 변경 | 504, latency | ConfigMap·Deployment 복구 |
| 잘못된 이미지 | `kubectl set image` | ImagePullBackOff | `kubectl rollout undo` |
| PDB 차단 | 2개 복제본 중 1개 NotReady 후 drain | eviction 거절 | 먼저 가용성 회복 |
| 용량 부족 | 큰 requests의 임시 Pod 생성 | Pending, FailedScheduling | Pod 삭제 또는 용량 증설 |

노드 실습의 자동화 버전은 `make resilience-kind`다. 갑작스러운 노드 정지는 감지 유예 시간이 길어서 수동 항목으로 둔다.

## 관찰

모든 장애에서 같은 순서를 사용한다.

```bash
kubectl get deployment,pod -n shop -o wide
kubectl describe pod -n shop POD_NAME
kubectl get events -n shop --sort-by=.lastTimestamp
kubectl logs -n shop POD_NAME --previous
```

판단은 사용자 영향 → 변경 이력 → 원하는 상태와 실제 상태 → Events → 로그·메트릭 순서로 좁힌다.

## 장애와 복구

복구는 원인 제거, 선언 상태 재적용, rollout 대기, API 확인까지 끝나야 한다. 증상별 명령은 [CrashLoop/Pending](../../docs/runbooks/crashloop-pending.md), [OOM/Eviction](../../docs/runbooks/oom-eviction.md), [노드 drain](../../docs/runbooks/node-drain.md)을 사용한다.

## 정리

```bash
kubectl apply -k deploy/overlays/kind
kubectl rollout status deployment -n shop --timeout=180s
kubectl get pdb -n shop
```
