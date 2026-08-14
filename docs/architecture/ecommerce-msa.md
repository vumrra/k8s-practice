# Ecommerce MSA 아키텍처

## 목적

도메인 완성품이 아니라 Kubernetes 서비스 디스커버리, 실패 전파, autoscaling, NetworkPolicy, observability와 service mesh를 실습하는 작은 호출 그래프다.

```text
Client
  -> gateway
       ├── GET /api/products -> catalog
       └── POST /api/orders  -> orders
                                  ├── inventory
                                  └── payments
```

| 서비스 | 책임 | 외부 호출 |
|---|---|---|
| gateway | 외부 API와 내부 route | catalog, orders |
| catalog | 고정 상품 조회 | 없음 |
| inventory | 결정론적 재고 조회 | 없음 |
| orders | 재고 확인 후 결제 조정 | inventory, payments |
| payments | 결제 승인 모사 | 없음 |

## 호출 계약

```bash
curl http://127.0.0.1:8080/api/products
curl -X POST http://127.0.0.1:8080/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"product_id":"pencil","quantity":1,"amount":1.5}'
```

orders는 dependency client timeout을 2초로 제한한다. inventory 5xx는 502, timeout은 504, 재고 부족은 409로 변환한다. 결제 POST는 idempotency 계약이 없으므로 자동 retry하지 않는다.

## Kubernetes 배치

- 서비스별 Deployment 2 replicas와 ClusterIP Service
- hostname topology spread와 PDB `minAvailable: 1`
- default deny NetworkPolicy 후 위 호출 그래프와 DNS만 허용
- 모든 workload non-root, read-only root filesystem, token mount 비활성화

데이터는 메모리 샘플이다. Pod 재시작 시 주문 상태가 사라지므로 실제 데이터 안정성을 주장하지 않는다. 실전에서는 서비스별 데이터 소유, transaction 경계, outbox/idempotency, schema migration과 backup을 별도로 설계한다.

## 실패 전파 확인

```bash
kubectl scale -n shop deployment/inventory --replicas=0
kubectl port-forward -n shop service/gateway 8080:8080
# 주문 요청은 gateway -> orders -> inventory에서 실패한다.
kubectl apply -k deploy/overlays/kind
```

실습: [Lab 13](../../labs/13-ecommerce-msa/README.md)
