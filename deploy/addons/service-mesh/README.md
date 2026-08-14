# Istio 선택 실습

서비스 메시가 실제로 필요한지 확인한 뒤에만 설치한다. 기본 MSA는 Kubernetes Service, NetworkPolicy, 애플리케이션 타임아웃만으로 동작한다.

## 설치

```bash
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update
helm upgrade --install istio-base istio/base \
  --namespace istio-system --create-namespace --wait
helm upgrade --install istiod istio/istiod \
  --namespace istio-system \
  --values deploy/addons/service-mesh/values.yaml --wait
kubectl label namespace shop istio-injection=enabled --overwrite
kubectl rollout restart deployment -n shop
kubectl rollout status deployment -n shop --timeout=180s
```

## 검증

```bash
kubectl get pods -n istio-system
kubectl get pods -n shop
kubectl get pod -n shop -l app=gateway \
  -o jsonpath='{.items[0].spec.containers[*].name}'
```

`istio-proxy`가 보여야 한다. mTLS와 트래픽 정책 실습은 [실습 14](../../../labs/14-service-mesh/README.md)를 따른다.

## 제거

```bash
kubectl label namespace shop istio-injection-
kubectl rollout restart deployment -n shop
helm uninstall istiod -n istio-system
helm uninstall istio-base -n istio-system
kubectl delete namespace istio-system
```

CRD는 다른 Istio 사용 여부를 확인하기 전에는 삭제하지 않는다.
