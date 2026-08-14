# Runbook — Backup, Restore와 Cluster 재구축

## 증상

object 삭제, data corruption, control plane/cluster 또는 storage 장애로 복구가 필요하다.

## 즉시 안전 확인

- 계속되는 write를 중단할지 판단하고 원본·snapshot을 덮어쓰지 않는다.
- incident 시간, 마지막 정상 시점, 요구 RPO/RTO를 기록한다.

## 진단

```bash
kubectl get all -A
kubectl get pv,pvc -A
kubectl get storageclass
kubectl get event -A --sort-by=.lastTimestamp
```

영향이 Kubernetes object, application data, PV, etcd, registry, DNS/identity 중 어디까지인지 나눈다. Git은 desired object를 복구하지만 database 내용을 복구하지 않는다.

## 복구

1. 새 cluster/control plane의 version과 add-on 호환성 확인
2. RKE2/kubeadm이면 검증된 etcd snapshot 절차로 control plane 복구
3. CNI, CSI, Gateway, DNS, identity와 registry 연결
4. Git/Kustomize/Argo CD에서 namespace와 application object 적용
5. storage backup을 새 PVC/DB에 restore
6. read-only 검증 후 트래픽 전환

제품별 restore 명령은 해당 배포판·storage 공식 절차를 사용한다. 운영 중인 etcd에 임의로 snapshot 파일을 덮어쓰지 않는다.

## 검증

object 개수보다 application read/write, data 시점, identity, DNS/TLS와 SLI를 확인한다. 실제 RTO/RPO를 기록한다.

## 예방

off-site/다른 account 보관, encryption, retention, access test, 정기 restore rehearsal와 backup failure alert를 둔다.
