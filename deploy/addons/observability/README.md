# kube-prometheus-stack

Prometheus, Alertmanager, Grafana, kube-state-metrics를 한 번에 설치한다. 제공 값은 노트북 실습용이며 운영 보존 기간·스토리지·인증·알림 수신자는 별도로 설계해야 한다.

## 설치

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace \
  --values deploy/addons/observability/values.yaml --wait
kubectl apply -f deploy/addons/observability/rules.yaml
```

저장소의 Grafana 비밀번호는 실습 전용이다. 실제 환경에서는 External Secrets 같은 비밀 관리 경계를 사용한다.

## 검증

```bash
kubectl get pods -n monitoring
kubectl get prometheusrule -n monitoring shop-practice
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:80
```

브라우저에서 `http://localhost:3000`을 열고 알림 규칙을 확인한다. 장애 분류는 [관측성 개념](../../../docs/concepts/10-observability.md), 대응 절차는 [runbook 목록](../../../docs/runbooks/crashloop-pending.md)에서 시작한다.

## 제거

```bash
kubectl delete -f deploy/addons/observability/rules.yaml --ignore-not-found
helm uninstall monitoring -n monitoring
kubectl delete namespace monitoring
```

차트 CRD는 Helm 제거 후에도 남을 수 있으므로 다른 사용자를 확인한 뒤 별도로 관리한다.
