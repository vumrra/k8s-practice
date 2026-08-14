# Runbook — Node Drain과 장애

## 증상

Node maintenance가 필요하거나 Node가 NotReady다.

## 즉시 안전 확인

- control plane/etcd quorum과 workload replicas를 확인한다.
- 한 번에 하나의 장애 도메인만 정비한다.

## 진단

```bash
kubectl get node
kubectl describe node NODE
kubectl get pod -A --field-selector spec.nodeName=NODE
kubectl get pdb -A
kubectl get event -A --sort-by=.lastTimestamp
```

DaemonSet, local storage, unmanaged Pod, PDB allowed disruptions와 다른 Node의 capacity를 확인한다.

## 복구 또는 정비

```bash
kubectl cordon NODE
kubectl drain NODE --ignore-daemonsets --delete-emptydir-data --timeout=10m
# OS/runtime/kubelet 정비
kubectl uncordon NODE
```

PDB가 막으면 `--disable-eviction`으로 우회하지 말고 replicas, unhealthy Pod, capacity와 PDB 의도를 먼저 수정한다. 갑작스러운 Node 장애는 drain과 달리 PDB가 막아주지 않는다.

## 검증

Node Ready, DaemonSet, CSI/CNI, CoreDNS와 실제 gateway 요청을 확인한다. Pod가 여러 Node에 다시 분산됐는지 본다.

## 예방

topology spread, capacity headroom, 정기 drain rehearsal, Node/lease alert와 maintenance 순서를 문서화한다.
