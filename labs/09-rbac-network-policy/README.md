# 09. RBAC과 NetworkPolicy

## 목표

API 권한을 제한하는 RBAC과 Pod 간 데이터 경로를 제한하는 NetworkPolicy의 경계를 구분한다.

## 준비

```bash
make deploy ENV=kind
kubectl get networkpolicy -n shop
```

## 실행

```bash
kubectl auth can-i get pods --as=system:serviceaccount:shop:gateway -n shop
kubectl auth can-i create deployments --as=system:serviceaccount:shop:gateway -n shop
kubectl run outside -n default --image=busybox:1.37 --restart=Never -- \
  sh -c 'wget -T 3 -q -O- http://gateway.shop:8080/api/products || true'
kubectl logs -n default outside
```

기본 ServiceAccount에는 토큰 자동 마운트를 끈 상태라 애플리케이션이 API를 호출할 수 없다. 외부 namespace의 Pod는 `shop` 기본 차단 정책 때문에 gateway로 들어갈 수 없다.

## 관찰

```bash
kubectl describe networkpolicy -n shop
kubectl get pod -n shop -l app.kubernetes.io/name=gateway \
  -o jsonpath='{.items[0].spec.automountServiceAccountToken}'
```

## 장애와 복구

정책 문제는 DNS, Service, EndpointSlice, NetworkPolicy 순으로 좁힌다. 임시로 전체 정책을 지우는 방식은 원인과 보안 경계를 동시에 없애므로 사용하지 않는다. [네트워크 runbook](../../docs/runbooks/network-dns.md)을 따른다.

## 정리

```bash
kubectl delete pod -n default outside --ignore-not-found
```
