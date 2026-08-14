# Runbook — Kubernetes Upgrade

## 증상

지원 종료 대응, 보안 patch 또는 기능 요구로 cluster와 add-on upgrade가 필요하다.

## 즉시 안전 확인

- 현재 version, 지원 기한, 변경 window와 rollback 불가 구간을 확인한다.
- etcd/application backup과 restore 검증 없이 시작하지 않는다.

## 진단과 사전 점검

```bash
kubectl version
kubectl get node
kubectl get --raw /metrics | grep apiserver_requested_deprecated_apis
kubectl get pdb -A
helm list -A
```

배포판 release note, version skew, deprecated API, CNI/CSI/Gateway/mesh/chart 호환성을 확인한다.

## 실행

1. staging/kind에서 manifest와 smoke test
2. control plane을 지원되는 한 minor씩 upgrade
3. core add-on, CNI, CSI와 controller
4. worker를 한 대씩 cordon/drain/upgrade/uncordon
5. kubectl, Helm과 운영 client

관리형 EKS와 RKE2는 각각 공식 절차를 따르고 upstream kubeadm 명령을 섞지 않는다.

## 검증

Node Ready 외에 DNS, Service, NetworkPolicy, PV attach, Gateway, metrics/alert와 실제 주문 요청을 확인한다. event와 SLI를 변경 전 기준과 비교한다.

## 예방

지원 종료 calendar, 정기 patch window, deprecation metric alert, upgrade rehearsal와 owner를 유지한다.
