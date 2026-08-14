# 16. AWS EKS 이식

## 목표

애플리케이션 base는 유지하고 이미지 저장소·Ingress·스토리지·IAM 경계만 EKS에 맞게 교체한다.

## 준비

AWS CLI, eksctl, 계정 권한과 예상 비용을 먼저 확인한다. EKS, EC2, ALB, NAT Gateway, EBS는 비용이 발생할 수 있으며 이 저장소는 AWS 생성 명령을 자동 실행하지 않는다.

## 실행

비용이 없는 정적 렌더부터 확인한다.

```bash
kubectl kustomize deploy/overlays/eks
kubectl kustomize deploy/overlays/eks | kubectl apply --dry-run=client -f -
```

실제 전환 시에만 다음 순서로 진행한다.

1. `clusters/eks/cluster.yaml`을 검토하고 사용자가 직접 `eksctl create cluster -f ...`를 실행한다.
2. 여섯 이미지를 ECR에 push하고 overlay의 image 이름을 고정 태그나 digest로 바꾼다.
3. AWS Load Balancer Controller와 Pod Identity를 설치한 뒤 ALB Ingress를 적용한다.
4. 영속 볼륨이 필요하면 EBS CSI 드라이버와 StorageClass를 선택한다.
5. rollout, API, 알림을 확인한 뒤 DNS를 점진적으로 전환한다.

## 관찰

```bash
kubectl get nodes -o wide
kubectl get deployment,pod -n shop
kubectl get ingress -n shop
kubectl get events -n shop --sort-by=.lastTimestamp
```

[로컬·EKS 비교](../../docs/architecture/local-eks-comparison.md)에서 그대로 유지되는 것과 교체되는 것을 확인한다.

## 장애와 복구

이미지 Pull 권한, ALB subnet/tag, Pod Identity association, EBS AZ 제약을 각각 분리해 진단한다. 배포 전 rollback 이미지와 DNS 복귀 조건을 정한다.

## 정리

실제 AWS 리소스를 만들었다면 `eksctl delete cluster` 전에 LoadBalancer와 PVC가 만든 외부 리소스를 확인한다. 삭제 여부는 비용과 데이터 보존 정책에 따라 사용자가 결정한다.
