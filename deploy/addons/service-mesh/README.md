# Istio 선택 실습

서비스 메시가 실제로 필요한지 확인한 뒤에만 설치한다. 기본 MSA는 Kubernetes Service, NetworkPolicy, 애플리케이션 타임아웃만으로 동작한다.

## 설치

```bash
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update
kubectl create namespace istio-system --dry-run=client -o yaml | kubectl apply -f -
kubectl label namespace istio-system \
  pod-security.kubernetes.io/enforce=privileged \
  pod-security.kubernetes.io/enforce-version=latest --overwrite
helm upgrade --install istio-base istio/base \
  --namespace istio-system --wait
helm upgrade --install istio-cni istio/cni \
  --namespace istio-system --wait
helm upgrade --install istiod istio/istiod \
  --namespace istio-system \
  --values deploy/addons/service-mesh/values.yaml --wait
kubectl apply -f - <<'YAML'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-istiod
  namespace: shop
spec:
  podSelector: {}
  policyTypes: [Egress]
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: istio-system
      ports:
        - {protocol: TCP, port: 15012}
YAML
kubectl label namespace shop istio-injection=enabled --overwrite
kubectl rollout restart deployment -n shop
kubectl rollout status deployment -n shop --timeout=180s
```

Istio CNI가 Node 네트워크를 설정하므로 `istio-system`만 privileged 정책을 사용한다. 애플리케이션 `shop` namespace는 restricted 정책을 유지한다.

## 검증

```bash
kubectl get pods -n istio-system
kubectl get pods -n shop
kubectl get pod -n shop -l app.kubernetes.io/name=gateway \
  -o jsonpath='{.items[0].spec.containers[*].name}'
```

`istio-proxy`가 보여야 한다. mTLS와 트래픽 정책 실습은 [실습 14](../../../labs/14-service-mesh/README.md)를 따른다.

## 제거

```bash
kubectl label namespace shop istio-injection-
kubectl rollout restart deployment -n shop
kubectl rollout status deployment -n shop --timeout=180s
kubectl delete networkpolicy -n shop allow-istiod --ignore-not-found
helm uninstall istiod -n istio-system
helm uninstall istio-cni -n istio-system
helm uninstall istio-base -n istio-system
kubectl delete namespace istio-system
```

CRD는 다른 Istio 사용 여부를 확인하기 전에는 삭제하지 않는다.
