# Kubernetes Practice

Pod부터 Go 이커머스 MSA, 장애 대응, HA/DR, 관측성, Service Mesh, GitOps, EKS 이식까지 단계적으로 연습하는 Kubernetes 모노레포다. base manifest는 공유하고 minikube, kind, EKS 차이만 overlay로 분리했다.

## 먼저 알아둘 점

Kubernetes는 로컬에서도 실행할 수 있다. 다만 macOS는 Linux kernel이 없으므로 Docker Desktop 같은 Linux VM 위에 Node를 만든다.

| 도구 | 정체 | 이 저장소에서의 용도 |
|---|---|---|
| Docker Desktop | macOS/Windows에서 Linux container와 VM을 제공하는 기반. 자체 Kubernetes 옵션도 있음 | image build와 minikube/kind 기반 |
| minikube | 로컬 Kubernetes cluster를 만들고 관리하는 도구 | 기초 학습, addon, 일상 개발 |
| kind | Docker container를 Kubernetes Node로 쓰는 도구 | 1 control-plane + 2 worker, drain/PDB, CI |
| kubectl | Kubernetes API client | 배포·관찰·진단 |
| Kustomize | base에 환경별 차이를 합성 | `kubectl`에 내장된 overlay 관리 |
| Helm | 외부 플랫폼 chart 관리자 | Traefik, Prometheus, Argo CD, Istio |

처음이면 minikube가 편하고, 다중 Node와 CI까지 보려면 kind가 적합하다. Docker Desktop 내장 Kubernetes를 써도 표준 manifest는 실행되지만 이 저장소의 cluster 생성·image load 명령은 minikube/kind를 기준으로 한다.

## 준비

- Docker Desktop 또는 Linux Docker Engine
- `kubectl`
- minikube 또는 kind
- Helm 4 또는 Helm 3
- Go 1.26.x. 로컬 Go가 없으면 `make test-docker` 사용

kind 전체 실습에는 Docker에 CPU 4개, 메모리 6GB 이상을 권장한다.

## 5분 시작: kind

```bash
make kind-up
make images
make load-images
make deploy
make smoke
```

다중 Node 복구도 확인하려면 다음을 실행한다. 스크립트는 현재 context가 `kind-k8s-practice`가 아니면 drain을 거부한다.

```bash
make resilience-kind
make clean-kind
```

## 5분 시작: minikube

```bash
make minikube-up
make images
make load-images ENV=minikube
make deploy ENV=minikube
make smoke
```

끝낸 뒤 cluster를 지운다.

```bash
minikube delete --profile k8s-practice
```

## 저장소 구조

```text
k8s-practice/
├── apps/                 Go 예제: hello-api, 이커머스 5개 서비스
├── internal/web/         공용 HTTP 응답·probe·graceful shutdown
├── build/                한 개의 재사용 Dockerfile
├── deploy/
│   ├── base/             환경과 무관한 표준 Kubernetes object
│   ├── overlays/         minikube, kind, EKS 차이
│   └── addons/           Gateway, metrics, observability, GitOps, mesh
├── clusters/             kind topology, 선택형 eksctl 설정
├── labs/                 01~16 실행형 실습
├── docs/
│   ├── concepts/         핵심 개념
│   ├── architecture/     MSA·신뢰성·환경 비교
│   ├── runbooks/         장애 대응 절차
│   └── onprem/           RKE2 운영 참고 구성
└── tests/                manifest, 문서, smoke, 복원력 검증
```

## 학습 순서

전체 완료 기준은 [학습 경로](docs/learning-path.md)에 있다.

| 단계 | 실습 | 핵심 |
|---|---|---|
| 기초 | [01 Pod](labs/01-pod/README.md) → [06 Storage](labs/06-storage-statefulset/README.md) | Pod, Deployment, Service, 설정, probe, resource, batch, StatefulSet |
| 활용 | [07 Gateway](labs/07-ingress-gateway/README.md) → [12 장애 복구](labs/12-failure-recovery/README.md) | 외부 진입, HPA, scheduling, RBAC, NetworkPolicy, 관측, rollout, PDB/drain |
| 실전 | [13 이커머스 MSA](labs/13-ecommerce-msa/README.md) → [16 EKS](labs/16-eks-migration/README.md) | 서비스 의존성, 선택형 mesh, GitOps, 클라우드 이식 |

처음부터 운영 관점을 같이 보려면 다음 문서를 우선 읽는다.

- [Pod 생명주기와 probe](docs/concepts/03-pod-lifecycle.md)
- [리소스와 스케줄링](docs/concepts/07-resources-scheduling.md)
- [보안과 멀티테넌시](docs/concepts/09-security-multitenancy.md)
- [관측성](docs/concepts/10-observability.md)
- [HA, 장애, 백업과 DR](docs/concepts/11-reliability-ha-dr.md)
- [유지보수와 업그레이드](docs/concepts/14-maintenance-upgrades.md)

## 예제 서비스

```text
client
  └─ gateway
      ├─ catalog
      └─ orders
          ├─ inventory
          └─ payments
```

모든 Deployment는 2개 복제본, readiness/liveness probe, requests/limits, non-root/read-only 보안 설정을 가진다. 이커머스에는 topology spread, PDB, default-deny와 호출 경로별 NetworkPolicy가 포함된다. 상세 실패 정책은 [이커머스 MSA 구조](docs/architecture/ecommerce-msa.md)를 참고한다.

## 자주 쓰는 명령

| 명령 | 동작 |
|---|---|
| `make test` | 로컬 Go 단위 테스트 |
| `make test-docker` | Go container에서 단위 테스트 |
| `make manifests` | 세 overlay와 HA·보안 규칙 정적 검증 |
| `make docs-check` | 필수 문서·실습·내부 링크 검증 |
| `make images` | 여섯 application image build |
| `make load-images [ENV=minikube]` | 로컬 cluster에 image 주입 |
| `make deploy [ENV=kind]` | overlay 적용 후 모든 Deployment 대기 |
| `make smoke` | 상품 조회와 주문 확정 API 확인 |
| `make resilience-kind` | Pod 삭제와 worker drain 복구 확인 |

상태를 볼 때는 다음 순서를 습관화한다.

```bash
kubectl get deployment,pod -A -o wide
kubectl describe pod -n NAMESPACE POD_NAME
kubectl get events -n NAMESPACE --sort-by=.lastTimestamp
kubectl logs -n NAMESPACE POD_NAME --previous
```

## EKS로 바꾸는 범위

애플리케이션 코드는 바꾸지 않는다. `deploy/base`를 유지하고 다음 경계만 교체한다.

- 로컬 image load → ECR push와 image rewrite
- port-forward/로컬 Gateway → AWS Load Balancer Controller와 ALB Ingress
- 로컬 StorageClass → EBS CSI
- 정적 AWS credential → EKS Pod Identity
- 로컬 Node → managed node group과 multi-AZ 용량 계획

먼저 비용 없이 렌더링한다.

```bash
kubectl kustomize deploy/overlays/eks
```

`clusters/eks/cluster.yaml`은 서울 리전의 private worker 3대, NAT Gateway, control-plane log를 포함한 참고 설정이다. EKS control plane, EC2, NAT Gateway, ALB, EBS, CloudWatch에는 비용이 발생할 수 있다. 저장소 명령은 cluster를 자동 생성하지 않으며 [실습 16](labs/16-eks-migration/README.md)에서 사용자가 검토 후 직접 실행한다.

## 온프레미스

운영 기준 예시는 RKE2 server 3대 + worker 3대 이상, 외부 API VIP, Cilium, MetalLB, Traefik, 기존 CSI 우선, Harbor, 관측·백업으로 정리했다. 서버·rack·전원·스토리지가 같은 장애 도메인이라면 복제본 수만 늘려도 HA가 되지 않는다. 설치 자동화는 환경별 네트워크·스토리지 차이가 크므로 포함하지 않고 [RKE2 운영 참고 구성](docs/onprem/rke2-production-reference.md)에 결정 항목만 담았다.

## 장애 대응 시작점

- Pod가 재시작하거나 Pending: [CrashLoop/Pending](docs/runbooks/crashloop-pending.md)
- 통신·DNS 실패: [Network/DNS](docs/runbooks/network-dns.md)
- OOM·Eviction: [OOM/Eviction](docs/runbooks/oom-eviction.md)
- 배포 실패: [Rollout/Rollback](docs/runbooks/rollout-rollback.md)
- 노드 유지보수: [Node drain](docs/runbooks/node-drain.md)
- 복구 목표와 데이터: [Backup/Restore](docs/runbooks/backup-restore.md)
- 버전 변경: [Upgrade](docs/runbooks/upgrade.md)

## GitHub 연결

빈 GitHub 저장소를 만든 다음 자신의 URL로 연결한다.

```bash
git remote add origin git@github.com:OWNER/k8s-practice.git
git push -u origin main
```

이미 `origin`이 있으면 `git remote -v`로 확인하고 `git remote set-url origin ...`을 사용한다. GitHub Actions는 단위·문서·manifest 검사와 3-node kind smoke/drain 복구를 수행한다.
