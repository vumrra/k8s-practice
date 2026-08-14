# 08. Storage와 Stateful Workload

## 학습 목표

- Volume, PV, PVC, StorageClass와 CSI의 관계를 설명한다.
- StatefulSet과 데이터 HA를 구분한다.

```text
Pod의 PVC 요청
 -> StorageClass와 CSI provisioner
 -> PersistentVolume 생성/선택
 -> Node에 attach/mount
 -> container volumeMount
```

- Volume: PodSpec 안의 mount 가능한 저장소 정의
- PVC: workload의 용량·access mode 요청
- PV: 클러스터 수준 저장소 자원
- StorageClass: provisioner, parameter, reclaim/volume binding 정책
- CSI: Kubernetes와 storage 구현의 표준 경계

## 주요 결정

- `ReadWriteOnce`는 일반적으로 한 Node에서 read-write이며 “Pod 하나만”이라는 뜻은 아니다.
- `WaitForFirstConsumer`는 Pod가 배치될 zone을 고려해 volume을 만든다.
- reclaim `Delete`는 PVC 삭제가 실제 volume 삭제로 이어질 수 있다.
- snapshot은 backup 전체 전략이 아니다. 다른 장애 도메인 보관과 restore 검증이 필요하다.

StatefulSet은 고정 ordinal, 순차 동작, 개별 PVC를 제공한다. 데이터 replication, leader election, quorum, consistency, backup을 자동 제공하지 않는다.

## 운영 확인

```bash
kubectl get storageclass
kubectl get pv,pvc -A
kubectl describe pvc -n practice data-writer-0
kubectl get volumeattachment
```

PVC Pending은 StorageClass 이름, provisioner 상태, topology, access mode와 quota를 확인한다. Node 장애 후 volume detach/attach에는 시간이 걸릴 수 있다.

## 체크리스트

- [ ] 데이터 수명주기를 Pod 삭제와 분리한다.
- [ ] reclaim policy를 확인한 후 PVC를 삭제한다.
- [ ] backup 성공이 아니라 restore 성공을 정기 검증한다.

실습: [Lab 06](../../labs/06-storage-statefulset/README.md)
