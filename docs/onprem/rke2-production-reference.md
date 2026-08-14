# 온프레미스 RKE2 운영 참고 구성

이 문서는 설치 자동화가 아니라 구성요소와 운영 책임을 정리한다. 실제 값은 서버, 네트워크, 스토리지, 보안과 지원 계약에 맞춰 설계한다.

## 권장 기준 구성

```text
운영자/kubectl
  -> L4 VIP 또는 기존 ADC: TCP 9345, 6443
      -> RKE2 server 3대: API, controller, scheduler, embedded etcd

사용자
  -> 사내 DNS
      -> MetalLB address
          -> Traefik Gateway
              -> worker 3대 이상의 application Pod

worker
  -> RKE2 agent + containerd
  -> Cilium CNI/NetworkPolicy/Hubble
  -> 기존 SAN/NAS CSI 또는 Longhorn
```

server는 odd number로 etcd quorum을 유지한다. 최소 3대 구성이 한 server 장애를 견디며, 고정 registration/API endpoint가 필요하다. server와 worker를 서로 다른 rack·전원·network 장애 도메인에 분산한다.

## 구성요소 선택

| 영역 | 기본 선택 | 운영 질문 |
|---|---|---|
| Kubernetes | RKE2 HA | 지원 version, patch/upgrade owner |
| Runtime | bundled containerd | registry mirror, disk 정리 |
| CNI | Cilium | kernel 호환, MTU, NetworkPolicy, BGP |
| Service IP | MetalLB | L2 또는 BGP, 할당 IP pool, router owner |
| HTTP entry | Traefik Gateway API | VIP, TLS, WAF, source IP |
| Storage | 기존 CSI 우선, 없으면 Longhorn | replication, failure domain, backup target |
| Registry | Harbor | HA, scanner, retention, air-gap sync |
| 인증서 | 사내 PKI/cert-manager | issuance, rotation, revocation |
| 관측 | Prometheus/Grafana/Alertmanager | 외부 장기 보관, alert route |
| GitOps | Argo CD | admin 권한, secret, disaster bootstrap |

회사가 이미 운영하는 L4, SAN/NAS와 PKI가 있으면 검증된 native 기능을 우선 사용한다. Kubernetes 안에 같은 기능을 중복 구축하지 않는다.

## 네트워크

- node, Pod, Service, MetalLB pool, 사내·VPN CIDR이 겹치지 않아야 한다.
- API VIP에서 모든 server의 6443/9345 health를 확인한다.
- MTU는 overlay/VLAN/VPN overhead를 반영하고 packet fragmentation을 시험한다.
- L2 MetalLB는 단순하지만 broadcast domain 제약이 있다. BGP는 network team과 ASN/route 정책이 필요하다.
- DNS, NTP와 private registry는 control plane을 포함한 모든 node에서 안정적으로 접근해야 한다.

## Storage와 데이터

StatefulSet이나 Longhorn 설치가 데이터 HA를 자동 완성하지 않는다. workload별 consistency, replica, quorum, failure domain, snapshot, off-site backup과 restore를 검토한다. etcd disk는 낮은 latency의 SSD를 사용하고 application volume과 경합을 피한다.

## 보안

- RKE2와 OS supported patch 유지
- node SSH·sudo 최소화와 중앙 identity
- API audit, Pod Security Admission, RBAC least privilege
- Secret encryption at rest와 외부 secret manager 검토
- private registry allowlist, image scan/signature와 immutable tag
- control plane, worker, storage, management network 분리
- air-gap이면 image, binary, checksum과 update 반입 절차 운영

## Backup과 DR

```text
매일/변경 전 etcd snapshot
 -> cluster 밖의 암호화된 위치로 복사
 -> retention과 backup job alert
 -> 분리 환경에서 정기 restore rehearsal

application data
 -> storage/database native backup
 -> 다른 장애 도메인 보관
 -> application read/write로 restore 검증
```

Git은 manifest를 복원하지만 etcd의 모든 runtime 상태나 PV data를 복원하지 않는다. Harbor image, DNS, PKI, storage credential과 GitOps bootstrap도 DR inventory에 포함한다.

## 운영 주기

- 매일: cluster/etcd/storage/registry/alert 상태
- 매주: capacity, certificate, failed backup, security event
- 매월: patch와 deprecated API, node drain rehearsal 일부
- 분기: backup restore, worker/rack 장애와 연락망 훈련
- upgrade 전후: compatibility, snapshot, canary node, smoke/SLI 비교

## 운영 준비 체크리스트

- [ ] control plane, worker, storage의 장애 도메인이 문서화됐다.
- [ ] API endpoint와 MetalLB IP의 소유자가 정해졌다.
- [ ] NTP/DNS/registry/storage 장애 시 연락과 우회 절차가 있다.
- [ ] PDB를 존중해 worker 하나를 drain할 capacity가 있다.
- [ ] etcd와 application data restore를 실제 수행했다.
- [ ] certificate, patch, vulnerability, capacity와 backup alert가 runbook에 연결된다.

참고: [RKE2 HA](https://docs.rke2.io/install/ha), [RKE2 Air-gap](https://docs.rke2.io/install/airgap), [MetalLB](https://metallb.io/), [Cilium](https://docs.cilium.io/en/stable/overview/intro/), [Longhorn](https://longhorn.io/docs/latest/what-is-longhorn/)
