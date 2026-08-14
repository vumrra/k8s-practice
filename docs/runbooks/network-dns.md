# Runbook — Service 통신과 DNS 실패

## 증상

connection refused, timeout, `no such host`, 502/504가 발생한다.

## 즉시 안전 확인

- 전체 경로인지 특정 service/namespace/node인지 범위를 나눈다.
- NetworkPolicy를 전체 삭제하기 전에 변경 시점과 영향 범위를 기록한다.

## 진단

```bash
kubectl get service,endpointslice,pod -n shop -o wide
kubectl describe service -n shop SERVICE
kubectl get networkpolicy -n shop
kubectl get pod -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=100
kubectl run -n shop netcheck --rm -it --restart=Never --image=curlimages/curl:8.17.0 -- sh
```

진단 Pod에서 `nslookup catalog`, `curl -v http://catalog:8080/products`를 실행한다. EndpointSlice가 비면 selector와 readiness, 연결 거부면 targetPort/process listen, timeout이면 NetworkPolicy/CNI/route를 본다.

## 복구

- label/selector, port, readiness 또는 DNS egress rule을 Git에서 수정해 적용
- 잘못된 NetworkPolicy revision rollback
- CoreDNS 자체 장애면 replicas, Node placement와 upstream resolver 확인

## 검증

gateway 상품 조회와 주문 경로를 호출하고 DNS query/error metric과 NetworkPolicy event를 확인한다.

## 예방

호출 그래프 기반 allow rule, DNS 별도 허용, endpoint와 synthetic API probe, 변경 전 렌더링 review를 유지한다.
