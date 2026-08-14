# 10. 메트릭·대시보드·알림

## 목표

증상 감지, 원인 탐색, runbook 연결을 하나의 운영 흐름으로 만든다.

## 준비

```bash
make deploy ENV=kind
```

## 실행

[관측성 애드온](../../deploy/addons/observability/README.md)을 설치한 뒤 gateway Pod를 반복 삭제해 경고 조건을 만든다.

```bash
kubectl delete pod -n shop -l app.kubernetes.io/name=gateway
kubectl rollout status deployment -n shop gateway --timeout=120s
kubectl get prometheusrule -n monitoring shop-practice -o yaml
```

## 관찰

```bash
kubectl port-forward -n monitoring svc/monitoring-kube-prometheus-prometheus 9090:9090
```

Prometheus의 `/alerts`와 `/targets`에서 규칙 상태와 수집 실패를 확인한다. 알림은 사용자 영향, 지속 시간, 담당자, runbook을 포함해야 한다.

## 장애와 복구

알림이 없으면 규칙보다 먼저 Prometheus target, ServiceMonitor label, 시계열 존재 여부를 확인한다. 알림 자체가 정상이어도 수집 경로가 끊기면 침묵할 수 있다.

## 정리

[관측성 애드온 제거 절차](../../deploy/addons/observability/README.md)를 따른다.
