# Runbook — CrashLoopBackOff와 Pending

## 증상

Pod가 Ready가 되지 않고 `CrashLoopBackOff`, `Pending`, `ImagePullBackOff`를 반복한다.

## 즉시 안전 확인

- 정상 replica가 남아 있는지, 사용자 오류율이 증가했는지 확인한다.
- 원인 확인 전에 전체 Deployment 재시작이나 강제 삭제를 반복하지 않는다.

## 진단

```bash
kubectl get pod -n NAMESPACE -o wide
kubectl describe pod -n NAMESPACE POD
kubectl logs -n NAMESPACE POD --all-containers
kubectl logs -n NAMESPACE POD --all-containers --previous
kubectl get event -n NAMESPACE --sort-by=.lastTimestamp
kubectl get node
kubectl get pvc -n NAMESPACE
```

CrashLoop은 exit code, signal, probe와 이전 로그를 본다. Pending은 scheduler event의 insufficient CPU/memory, taint, affinity, topology, PVC를 본다. ImagePullBackOff는 이름·tag·registry 인증·architecture를 확인한다.

## 복구

- bad image/config rollout: `kubectl rollout undo -n NAMESPACE deployment/NAME`
- capacity: 불필요한 workload를 줄이거나 Node 추가 후 재스케줄
- probe: 실제 endpoint와 startup 시간을 고쳐 새 revision 배포
- PVC: StorageClass/provisioner/topology를 고친다. 데이터 확인 없이 PVC를 삭제하지 않는다.

## 검증

```bash
kubectl rollout status -n NAMESPACE deployment/NAME
kubectl get pod -n NAMESPACE
kubectl get event -n NAMESPACE --sort-by=.lastTimestamp
```

실제 사용자 API도 호출한다.

## 예방

CI image build, manifest 렌더링, realistic requests, startup/readiness 분리, immutable tag, capacity alert를 둔다.
