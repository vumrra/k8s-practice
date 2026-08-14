# Traefik + Gateway API

서비스를 클러스터 밖에 노출하고 `GatewayClass`, `Gateway`, `HTTPRoute`의 역할을 실습한다. 로컬 학습용 NodePort 설정이며 운영 환경에서는 클라우드 LoadBalancer나 온프레미스 로드밸런서 정책에 맞춘다.

## 설치

```bash
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  --namespace traefik --create-namespace \
  --values deploy/addons/gateway/values.yaml --wait
```

차트 버전은 실습 재현성을 위해 팀에서 검증한 버전으로 고정하는 것이 좋다.

## 검증

```bash
kubectl get gatewayclass
kubectl get pods -n traefik
kubectl get svc -n traefik
```

라우팅 생성과 요청 확인은 [실습 07](../../../labs/07-ingress-gateway/README.md)에서 진행한다.

## 제거

```bash
helm uninstall traefik -n traefik
kubectl delete namespace traefik
```

Gateway API CRD는 다른 컨트롤러도 사용할 수 있으므로 자동 삭제하지 않는다.
