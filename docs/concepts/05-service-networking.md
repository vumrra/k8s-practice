# 05. Service, DNS와 네트워크

## 학습 목표

- Pod IP, Service, EndpointSlice, DNS와 Gateway의 책임을 구분한다.
- NetworkPolicy가 허용 목록 방식임을 이해한다.

## 요청 흐름

```text
외부 Client -> LoadBalancer/NodePort -> Gateway
Gateway -> ClusterIP Service DNS -> Ready Pod EndpointSlice
Pod -> CNI network -> 다른 Pod
```

Service는 selector에 맞는 Ready Pod를 EndpointSlice로 묶고 안정적인 가상 IP와 DNS를 제공한다. 같은 namespace에서는 `catalog`, 다른 namespace에서는 `catalog.shop.svc.cluster.local`로 접근한다.

## Service 유형

- ClusterIP: 클러스터 내부 기본값
- NodePort: 모든 Node의 정해진 port로 노출
- LoadBalancer: cloud controller나 MetalLB 같은 구현이 외부 주소 제공
- ExternalName: DNS CNAME 반환, proxy가 아님
- Headless(`clusterIP: None`): 개별 Pod DNS가 필요한 StatefulSet 등에 사용

Ingress는 HTTP routing API이고 실제 동작에는 controller가 필요하다. Gateway API는 infra 담당자의 Gateway와 app 담당자의 HTTPRoute 책임을 더 명확히 분리한다.

## NetworkPolicy

policy를 선택하지 않은 Pod는 기본 허용이다. default deny를 적용한 후 DNS와 필요한 서비스 경로를 명시적으로 연다. CNI가 NetworkPolicy를 지원하지 않으면 object가 있어도 차단되지 않는다.

```bash
kubectl get service,endpointslice -n shop
kubectl run -n shop netcheck --rm -it --restart=Never --image=curlimages/curl:8.17.0 -- curl -sS http://catalog:8080/products
kubectl get networkpolicy -n shop
```

## 흔한 장애

- Service selector와 Pod label 불일치: endpoint 없음
- targetPort와 containerPort 불일치: 연결 거부
- readiness 실패: endpoint에서 제거
- DNS egress 차단: Service 이름 해석 실패
- 잘못된 NetworkPolicy source selector: timeout

## 체크리스트

- [ ] 앱 간 호출은 Pod IP가 아닌 Service DNS를 사용한다.
- [ ] LoadBalancer object만 만들면 온프레미스 외부 IP가 생기지 않음을 안다.
- [ ] default deny 뒤 DNS와 호출 그래프를 허용한다.

실습: [Lab 02](../../labs/02-deployment-service/README.md), [Lab 07](../../labs/07-ingress-gateway/README.md), [Lab 09](../../labs/09-rbac-network-policy/README.md)
