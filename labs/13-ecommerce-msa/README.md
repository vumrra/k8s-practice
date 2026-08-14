# 13. 이커머스 MSA 통합

## 목표

gateway → orders → inventory/payments 호출과 읽기 경로 gateway → catalog를 실제 요청으로 확인한다.

## 준비

```bash
make images
make load-images
make deploy ENV=kind
kubectl rollout status deployment -n shop --timeout=180s
```

## 실행

```bash
kubectl port-forward -n shop svc/gateway 18080:8080
```

다른 터미널에서 실행한다.

```bash
curl --fail http://localhost:18080/api/products
curl --fail -X POST http://localhost:18080/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"product_id":"pencil","quantity":2,"amount":3}'
```

## 관찰

```bash
kubectl logs -n shop -l app.kubernetes.io/name=gateway --tail=20
kubectl logs -n shop -l app.kubernetes.io/name=orders --tail=20
kubectl get endpointslice -n shop
```

호출 구조와 실패 의미는 [이커머스 아키텍처](../../docs/architecture/ecommerce-msa.md)에 정리되어 있다.

## 장애와 복구

inventory를 0개로 줄이고 주문을 보내면 timeout/실패가 gateway까지 어떻게 전달되는지 확인한다. payments 요청은 중복 결제 위험 때문에 자동 재시도하지 않는다.

```bash
kubectl scale deployment -n shop inventory --replicas=0
kubectl scale deployment -n shop inventory --replicas=2
kubectl rollout status deployment -n shop inventory --timeout=120s
```

## 정리

port-forward를 `Ctrl-C`로 종료한다.
