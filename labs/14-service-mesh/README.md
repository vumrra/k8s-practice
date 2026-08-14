# 14. 선택형 서비스 메시

## 목표

Istio sidecar 주입과 namespace 단위 strict mTLS를 확인하고 메시가 추가하는 비용을 판단한다.

## 준비

[Istio 애드온](../../deploy/addons/service-mesh/README.md)을 설치한다. 노트북 자원이 부족하면 이 실습은 건너뛴다.

## 실행

```bash
kubectl apply -f - <<'YAML'
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: strict
  namespace: shop
spec:
  mtls:
    mode: STRICT
YAML
kubectl rollout status deployment -n shop --timeout=180s
```

## 관찰

```bash
kubectl get peerauthentication -n shop
kubectl get pods -n shop -o jsonpath='{range .items[*]}{.metadata.name}{"  "}{.spec.containers[*].name}{"\n"}{end}'
kubectl logs -n shop -l app=gateway -c istio-proxy --tail=20
```

프록시 CPU·메모리, 장애 지점, 업그레이드 책임이 늘어난다. 필요한 기능이 mTLS 하나뿐이라면 CNI·플랫폼 기능과 비교한다.

## 장애와 복구

sidecar가 없는 기존 Pod는 strict mTLS 통신에 실패한다. namespace label 후 모든 Deployment가 재시작됐는지 확인한다.

## 정리

```bash
kubectl delete peerauthentication -n shop strict
```

전체 제거는 [애드온 제거 절차](../../deploy/addons/service-mesh/README.md)를 따른다.
