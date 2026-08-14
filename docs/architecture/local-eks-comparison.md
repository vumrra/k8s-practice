# minikube, kind와 EKS 경계

## 공통으로 유지하는 것

Deployment, ClusterIP Service, ConfigMap 참조, probe, resources, HPA, PDB, security context, RBAC와 NetworkPolicy는 Kubernetes 표준 object이므로 `deploy/base`에서 공유한다.

| 영역 | minikube | kind | EKS |
|---|---|---|---|
| 목적 | 개인 학습·개발 | 다중 노드·CI | AWS 운영 |
| Node | 기본 단일 로컬 Node | Docker container Node | EC2 managed node/Fargate |
| 이미지 | `minikube image load` | `kind load docker-image` | ECR push |
| 외부 트래픽 | port-forward/tunnel | port-forward/port mapping | AWS Load Balancer Controller |
| StorageClass | local provisioner | local-path 계열 | EBS CSI 또는 다른 CSI |
| workload AWS 권한 | 없음 | 없음 | EKS Pod Identity 권장 |
| HA | laptop가 단일 장애점 | Docker host가 단일 장애점 | multi-AZ와 node 배치 설계 |

## EKS에서 바뀌는 부분

```text
image name -> ACCOUNT.dkr.ecr.REGION.amazonaws.com/repository:tag
Gateway/Ingress -> ALB 또는 NLB/Gateway API controller
StorageClass -> EBS CSI provisioner와 zone
ServiceAccount -> Pod Identity association
DNS/TLS -> Route 53, ACM 또는 cert-manager 선택
Node capacity -> managed node group, autoscaler/Karpenter 선택
```

AWS resource를 만드는 controller는 IAM 권한과 subnet/security group tag를 필요로 한다. 단순히 local YAML을 적용한다고 ALB, EBS와 IAM이 자동 준비되지는 않는다.

## 이 저장소의 원칙

```bash
kubectl kustomize deploy/overlays/minikube
kubectl kustomize deploy/overlays/kind
kubectl kustomize deploy/overlays/eks
```

세 명령은 credential이나 비용 없이 렌더링된다. 실제 EKS 생성과 ECR push는 [Lab 16](../../labs/16-eks-migration/README.md)에서 사용자가 명시적으로 실행한다.
