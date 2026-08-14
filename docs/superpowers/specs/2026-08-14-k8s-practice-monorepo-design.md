# Kubernetes Practice Monorepo Design

## 목표

`k8s-practice`는 Kubernetes를 처음 접하는 사람이 Pod부터 시작해 로컬 다중 노드, 이커머스 MSA, 관측, 보안, GitOps, 서비스 메시, EKS 이식까지 순서대로 실습할 수 있는 Go 기반 모노레포다. 모든 실습은 복사된 예제 묶음이 아니라 공통 애플리케이션과 공통 Kubernetes 리소스를 재사용한다.

## 범위

- 로컬 기본 환경: Docker Desktop 위의 minikube
- 로컬 다중 노드 및 CI: kind
- 클라우드 이전 예제: Amazon EKS용 Kustomize overlay와 최소 `eksctl` 설정
- 온프레미스: RKE2 운영 구성 참고 문서만 제공
- 애플리케이션: Go 표준 라이브러리 기반 Hello API와 이커머스 MSA
- 패키징: 애플리케이션은 Kustomize, 외부 애드온은 Helm
- 자동 검증: Go 테스트, Kustomize 렌더링, kind 배포, HTTP 스모크 테스트

프론트엔드, 실제 결제 연동, 완전한 상품·주문 데이터베이스, Terraform 기반 AWS 인프라, 온프레미스 자동 설치, 멀티클라우드 추상화, 커스텀 Kubernetes Operator는 포함하지 않는다. 이들은 Kubernetes 학습 목표에 비해 코드와 운영 부담이 크다.

## 핵심 설계 결정

### 하나의 점진적 모노레포

단계별로 프로젝트를 복제하지 않는다. `apps/`의 동일한 실행 파일과 `deploy/base/`의 동일한 리소스를 `labs/`에서 점진적으로 확장한다. 이 방식은 YAML 중복과 단계 간 설정 불일치를 줄이고 실제 운영 저장소의 변경 흐름을 보여준다.

### Kustomize와 Helm의 역할 분리

- 직접 만든 Go 애플리케이션: `deploy/base`와 `deploy/overlays`의 Kustomize 사용
- Metrics Server, Traefik, Prometheus, Grafana, Argo CD, 서비스 메시 같은 외부 제품: 공식 Helm chart 사용

Kustomize는 `kubectl`에 포함된 기능을 사용한다. 애플리케이션용 Helm chart를 별도로 만들지 않아 같은 설정을 두 번 관리하지 않는다.

### 환경 이식성

Deployment, ClusterIP Service, ConfigMap, Secret 참조, Probe, 리소스 제한, HPA, PDB, RBAC, NetworkPolicy는 `deploy/base`에 둔다. 이미지 주소, 외부 진입점, StorageClass와 클라우드별 annotation만 overlay에서 바꾼다.

```text
deploy/base
    ├── deploy/overlays/minikube
    ├── deploy/overlays/kind
    └── deploy/overlays/eks
```

## 디렉터리 구조

```text
k8s-practice/
├── .github/workflows/verify.yaml
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── apps/
│   ├── hello-api/
│   └── ecommerce/
│       ├── gateway/
│       ├── catalog/
│       ├── inventory/
│       ├── orders/
│       └── payments/
├── clusters/
│   ├── kind/config.yaml
│   └── eks/cluster.yaml
├── deploy/
│   ├── base/
│   │   ├── hello/
│   │   └── ecommerce/
│   ├── overlays/
│   │   ├── minikube/
│   │   ├── kind/
│   │   └── eks/
│   └── addons/
│       ├── gateway/
│       ├── metrics/
│       ├── observability/
│       ├── gitops/
│       └── service-mesh/
├── labs/
│   ├── 01-pod/
│   ├── 02-deployment-service/
│   ├── 03-config-secret/
│   ├── 04-probe-resource/
│   ├── 05-job-cronjob/
│   ├── 06-storage-statefulset/
│   ├── 07-ingress-gateway/
│   ├── 08-autoscaling-scheduling/
│   ├── 09-rbac-network-policy/
│   ├── 10-observability/
│   ├── 11-rollout-rollback/
│   ├── 12-failure-recovery/
│   ├── 13-ecommerce-msa/
│   ├── 14-service-mesh/
│   ├── 15-gitops/
│   └── 16-eks-migration/
├── docs/
│   ├── learning-path.md
│   ├── concepts/
│   ├── architecture/
│   ├── runbooks/
│   ├── onprem/
│   └── superpowers/
└── tests/smoke.sh
```

## 애플리케이션 구성

### Hello API

Hello API는 Pod, Deployment, Service, ConfigMap, Secret, Probe, 리소스 제한과 HPA를 학습하기 위한 최소 HTTP 서버다.

- `GET /`: 애플리케이션 이름, 버전, Pod 이름 반환
- `GET /healthz`: 프로세스 생존 상태 반환
- `GET /readyz`: 요청 처리 가능 상태 반환
- `GET /config`: 공개 가능한 환경 설정만 반환

서버는 HTTP 타임아웃과 종료 신호를 처리하고, 종료 시 새 요청을 받지 않은 뒤 진행 중인 요청을 정리한다.

### 이커머스 MSA

각 서비스는 독립 Deployment와 ClusterIP Service로 실행된다. 도메인 로직은 Kubernetes 통신과 장애 동작을 관찰할 수 있는 수준으로만 구현한다.

- `gateway`: 외부 요청의 단일 진입점이며 내부 서비스로 라우팅
- `catalog`: 고정된 상품 목록과 상품 조회 제공
- `inventory`: 상품별 재고 확인 응답 제공
- `orders`: 주문 요청을 검증하고 inventory와 payments 호출을 조정
- `payments`: 성공·거절을 결정론적으로 모사

기본 주문 흐름은 다음과 같다.

```text
Client -> Gateway -> Orders -> Inventory
                         └──> Payments
Client <- Gateway <- Orders
```

서비스 간 요청에는 명시적인 타임아웃을 적용하고, 하위 서비스 오류는 일관된 HTTP 상태와 짧은 오류 본문으로 변환한다. 데이터는 메모리에 유지해 데이터베이스 운영이 Kubernetes 학습을 가리지 않도록 한다. Pod 재시작으로 메모리 상태가 사라지는 현상은 상태 저장 학습에서 의도적으로 다룬다.

## 학습 단계

1. Pod와 컨테이너 실행 단위
2. Deployment, ReplicaSet, Service와 서비스 디스커버리
3. ConfigMap, Secret과 설정 주입
4. Liveness, Readiness, Startup Probe와 리소스 QoS
5. Job과 CronJob
6. Volume, PV, PVC와 StatefulSet
7. Traefik과 Gateway API를 사용한 외부 라우팅
8. HPA, affinity, taint, toleration, topology spread
9. ServiceAccount, RBAC, Pod Security, NetworkPolicy
10. Prometheus와 Grafana를 사용한 메트릭 관측
11. RollingUpdate, rollback, blue-green과 수동 canary
12. CrashLoopBackOff, Pending, 서비스 장애, 노드 drain 복구
13. 이커머스 MSA 배포와 서비스 간 통신
14. 서비스 메시의 트래픽 제어, mTLS와 관측을 선택 실습
15. Argo CD를 사용한 GitOps 동기화
16. 공통 base를 EKS overlay로 배포하는 이전 실습

각 lab은 목표, 사전 조건, 실행 명령, 관찰할 상태, 고장 유도, 복구 명령, 정리 명령을 포함한다. 답을 숨기는 별도 starter/solution 복제본은 만들지 않는다.

## 문서 구성

### 기본 개념

`docs/concepts`는 다음 내용을 독립 문서로 설명한다.

- 클러스터, control plane, worker node, kubelet, container runtime
- Pod, Deployment, ReplicaSet, DaemonSet, StatefulSet
- Service, EndpointSlice, CoreDNS, Ingress, Gateway API
- ConfigMap, Secret, 환경 변수와 volume mount
- Probe, requests, limits, QoS와 OOM
- Job, CronJob
- Volume, PV, PVC, StorageClass
- Scheduler, affinity, taint, toleration, HPA, PDB
- ServiceAccount, RBAC, Pod Security와 NetworkPolicy
- 로그, 메트릭, 트레이스와 관측 가능성
- 서비스 메시의 sidecar/ambient 개념, mTLS, retry, timeout, traffic split과 적용 기준
- Helm, Kustomize와 각각의 사용 범위
- GitOps와 Argo CD 동기화 모델

### 아키텍처와 운영 문서

- `docs/architecture/ecommerce-msa.md`: 서비스 책임, 호출 흐름, 실패 전파
- `docs/architecture/local-eks-comparison.md`: 로컬과 EKS에서 공통인 부분과 달라지는 부분
- `docs/runbooks`: CrashLoopBackOff, Pending Pod, 통신 실패, rollback, node drain, 백업·복구 점검 순서
- `docs/onprem/rke2-production-reference.md`: RKE2 HA, API load balancer, CNI, MetalLB, Gateway, 스토리지, 레지스트리, 모니터링과 백업의 참고 구성

온프레미스 문서는 설계 참고 자료이며 서버 프로비저닝 자동화나 설치 스크립트를 포함하지 않는다.

## 배포와 애드온

기본 배포는 외부 애드온 없이 동작한다. 각 애드온은 해당 lab에서만 설치한다.

- Gateway: Traefik과 Kubernetes Gateway API
- Metrics: Metrics Server
- Observability: Prometheus와 Grafana
- GitOps: Argo CD
- Service mesh: 기본 경로에서 제외된 Istio 선택 실습

서비스 메시 실습은 서비스 간 호출이 실제로 존재하는 이커머스 MSA에만 적용한다. 단순 Hello API에는 mesh를 추가하지 않는다.

## 구현 단위

전체 범위를 한 번에 결합하지 않고 다음 네 단위로 구현한다. 각 단위는 자체 검증 명령과 실행 가능한 결과를 남긴다.

1. Foundation: 저장소 도구, Hello API, minikube, 기본 개념 문서와 lab 01~06
2. Platform: Gateway, scheduling, security, observability와 장애 대응 lab 07~12
3. Ecommerce: 다섯 Go 서비스, MSA 배포, Istio 선택 실습과 lab 13~14
4. Delivery: kind CI, Argo CD, EKS overlay, 이전 문서와 lab 15~16

온프레미스 RKE2 참고 문서는 Foundation에 포함하지만 실제 클러스터 자동화와 배포 검증에서는 제외한다.

## 보안 기본값

- 컨테이너는 non-root 사용자로 실행
- 가능한 컨테이너는 read-only root filesystem 사용
- Linux capability를 기본 제거
- CPU와 memory requests/limits 지정
- Secret 값을 Git에 커밋하지 않고 예제 키와 생성 명령만 제공
- ServiceAccount와 RBAC는 필요한 권한만 부여
- MSA namespace는 기본 deny NetworkPolicy에서 필요한 통신만 허용

## 검증 전략

- `go test ./...`: HTTP handler, timeout과 종료 동작 검증
- `make manifests`: 모든 Kustomize overlay 렌더링 검증
- `make build`: 모든 Go 바이너리와 컨테이너 이미지 빌드
- `make smoke-kind`: kind 클러스터 생성, 이미지 로드, 배포, readiness 대기, HTTP 호출 검증
- GitHub Actions: Go 테스트, manifest 렌더링, kind 스모크 테스트 실행

스모크 테스트는 Hello API 응답과 이커머스 상품 조회·주문 흐름을 확인한다. CI 종료 시 생성한 kind 클러스터를 제거한다.

## 성공 기준

- 새 사용자가 README와 learning path만 따라 minikube에 Hello API를 배포할 수 있다.
- 각 lab은 독립적으로 시작 조건과 정리 절차를 제공한다.
- 동일한 이커머스 base가 minikube와 kind에서 동작한다.
- kind 기반 CI가 Go 테스트, 렌더링과 핵심 HTTP 흐름을 검증한다.
- EKS에서 변경해야 하는 이미지, ingress/load balancer, storage와 identity 경계를 문서와 overlay로 확인할 수 있다.
- Pod부터 서비스 메시까지 핵심 개념 문서가 실습 파일을 직접 가리킨다.
- 온프레미스 RKE2 구성은 구현 없이 운영 구성요소와 책임을 설명한다.
