# 07. Gateway API로 외부 요청 받기

## 목표

Traefik을 Gateway API 구현체로 설치하고 `Gateway`와 `HTTPRoute`로 gateway 서비스를 노출한다.

## 준비

```bash
make deploy ENV=kind
```

[Gateway 애드온](../../deploy/addons/gateway/README.md)을 먼저 설치한다.

## 실행

```bash
kubectl apply -f - <<'YAML'
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: shop
  namespace: shop
spec:
  gatewayClassName: traefik
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: Same
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: shop
  namespace: shop
spec:
  parentRefs:
    - name: shop
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api
      backendRefs:
        - name: gateway
          port: 8080
YAML
kubectl port-forward -n traefik svc/traefik 18080:80
```

다른 터미널에서 확인한다.

```bash
curl http://localhost:18080/api/products
```

## 관찰

```bash
kubectl describe gateway -n shop shop
kubectl describe httproute -n shop shop
```

`Accepted=True`, `ResolvedRefs=True`가 아니면 Events와 참조 이름·포트를 확인한다.

## 장애와 복구

HTTPRoute의 서비스 포트를 `9999`로 바꾸고 `ResolvedRefs` 상태와 503 응답을 관찰한 뒤 `8080`으로 복구한다.

## 정리

```bash
kubectl delete httproute,gateway -n shop shop
```
