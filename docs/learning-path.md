# Kubernetes 학습 경로

이 저장소는 YAML 암기보다 “원하는 상태를 선언하고, 관찰하고, 실패시키고, 복구하는” 흐름을 반복한다.

## 준비 도구

| 도구 | 역할 |
|---|---|
| Docker Desktop | macOS에서 Linux 컨테이너를 실행하는 기반 |
| minikube | 기본 로컬 Kubernetes 실습 |
| kind | 다중 노드와 CI 실습 |
| kubectl | Kubernetes API 조작 |
| Kustomize | 환경별 manifest 차이 관리, kubectl 내장 |
| Helm | 외부 애드온 설치 |

## 권장 순서

| 단계 | Lab | 결과 |
|---|---|---|
| 기초 | 01~06 | Pod, controller, Service, 설정, batch, storage 이해 |
| 플랫폼 | 07~12 | Gateway, autoscaling, 보안, 관측, 배포, 장애 대응 |
| 실전 | 13~16 | MSA, service mesh, GitOps, EKS 이전 |

각 lab은 앞 단계를 재사용한다. 문제가 생기면 바로 답을 재적용하기 전에 `kubectl get`, `describe`, `logs`, `events`로 현재 상태와 원인을 확인한다.

## 완료 기준

- `spec`과 `status`의 차이를 설명할 수 있다.
- Pod가 아닌 controller로 애플리케이션을 운영하는 이유를 안다.
- Service DNS, Gateway, NetworkPolicy의 책임을 구분한다.
- probe, requests/limits, HPA, PDB와 topology spread의 한계를 안다.
- HA, backup과 DR 및 SLI/SLO와 RTO/RPO를 구분한다.
- 실패한 rollout, node drain, OOM, DNS와 dependency 장애를 진단한다.
- 같은 base를 minikube, kind, EKS overlay로 렌더링한다.

첫 실습: [Lab 01 — Pod](../labs/01-pod/README.md)
