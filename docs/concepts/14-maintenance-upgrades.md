# 14. 유지보수와 Upgrade

## 학습 목표

- 지원 버전, version skew, deprecated API를 점검한다.
- control plane과 worker upgrade를 안전하게 수행하는 순서를 안다.

Kubernetes minor는 영구 지원되지 않는다. 지원 중인 minor의 최신 patch를 사용하고 배포판·cloud provider의 지원 표를 우선한다. `kubectl`, kubelet, control plane과 Helm도 각자의 version skew 정책이 있다.

## 사전 점검

1. release note와 API deprecation 확인
2. add-on, CNI, CSI, ingress/gateway, mesh 호환성 확인
3. etcd와 application data backup 및 restore 가능성 확인
4. PDB, capacity, maintenance window와 rollback 기준 확인
5. staging/kind에서 manifest·smoke test

```bash
kubectl version
kubectl get --raw /metrics | grep apiserver_requested_deprecated_apis
kubectl get pdb -A
kubectl get node
```

## 일반 순서

```text
control plane 한 minor씩
 -> core add-on/CNI/CSI
 -> worker를 한 대씩 cordon/drain/upgrade/uncordon
 -> kubectl과 운영 도구
 -> smoke, SLI, event 확인
```

```bash
kubectl cordon NODE
kubectl drain NODE --ignore-daemonsets --delete-emptydir-data --timeout=10m
kubectl uncordon NODE
```

직접 구축한 cluster는 API server/etcd 인증서 만료와 rotation도 운영 항목이다. EKS와 RKE2는 각 제품 절차를 따르며 upstream 명령을 섞지 않는다.

## 실패 대비

- PDB가 drain을 막으면 강제 삭제보다 replicas/capacity/정책을 먼저 조정
- deprecated API는 control plane upgrade 전에 manifest와 controller부터 교체
- rollback이 지원되지 않는 control plane/etcd 변경을 가정하고 backup restore 절차 준비
- upgrade 후 Node Ready만 보지 말고 실제 요청과 DNS/storage/network 확인

## 체크리스트

- [ ] 한 번에 여러 minor를 건너뛰지 않는다.
- [ ] 지원 종료 전에 계획한다.
- [ ] backup 파일 존재가 아니라 restore를 검증한다.

실습: [Lab 12](../../labs/12-failure-recovery/README.md) · [공식 Upgrade 개요](https://kubernetes.io/docs/tasks/administer-cluster/cluster-upgrade/)
