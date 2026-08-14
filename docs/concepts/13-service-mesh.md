# 13. Service Mesh

## 학습 목표

- mesh의 control plane과 data plane을 구분한다.
- mTLS, traffic policy, telemetry의 이점과 비용을 판단한다.

Service Mesh는 애플리케이션 간 통신에 공통 정책을 적용하는 계층이다.

```text
Service A -> data plane proxy/ambient tunnel -> Service B
                    ^
              mesh control plane
```

sidecar mode는 각 Pod에 proxy를 넣고, ambient mode는 node/namespace 계층의 공유 dataplane을 사용한다. 제품과 mode에 따라 지원 기능과 자원 비용이 다르다.

## 활용

- workload identity 기반 mTLS
- timeout, retry, connection pool과 outlier detection
- weighted traffic split을 이용한 canary
- 서비스 간 metric과 trace context
- 서비스별 authorization policy

## 주의점

- mesh는 복제·DB HA·backup을 제공하지 않는다.
- retry는 부하를 증폭할 수 있고 payment 같은 non-idempotent 요청을 중복 처리할 수 있다.
- 앱 timeout보다 proxy timeout이 짧은지, 전체 deadline budget이 맞는지 확인한다.
- sidecar/tunnel 장애와 잘못된 정책이라는 새로운 장애 유형이 생긴다.
- 작은 시스템에서는 앱 표준 HTTP client와 NetworkPolicy만으로 충분할 수 있다.

## 도입 기준

여러 팀·언어에서 mTLS/traffic policy/telemetry를 일관되게 강제해야 하고 그 운영 비용을 맡을 platform team이 있을 때 검토한다. 서비스 수가 적고 요구가 단순하면 먼저 앱 timeout, structured logs, metrics와 NetworkPolicy를 정리한다.

이 저장소에서는 Istio를 `shop` namespace에만 선택 적용한다. 기본 MSA와 smoke test는 mesh 없이도 동작한다.

## 체크리스트

- [ ] 해결하려는 구체적 문제와 성공 metric이 있다.
- [ ] retry 대상 요청의 idempotency를 확인했다.
- [ ] mesh 장애 시 진단·우회·upgrade runbook이 있다.

실습: [Lab 14](../../labs/14-service-mesh/README.md) · [Istio 공식 문서](https://istio.io/latest/docs/)
