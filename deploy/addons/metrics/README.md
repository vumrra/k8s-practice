# Metrics Server

Metrics Server는 `kubectl top`과 HPA가 쓰는 CPU·메모리 사용량 API를 제공한다. 장기 보관이나 알림 시스템은 아니며 그 역할은 Prometheus가 담당한다.

## kind 설치

```bash
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm repo update
helm upgrade --install metrics-server metrics-server/metrics-server \
  --namespace kube-system \
  --set 'args={--kubelet-insecure-tls}' --wait
```

kind의 자체 서명 kubelet 인증서를 위한 실습 전용 옵션이다. 운영에서는 이 옵션을 쓰지 말고 kubelet serving 인증서를 신뢰할 CA와 `--kubelet-certificate-authority`를 구성한다.

## minikube 설치

minikube에서는 내장 애드온을 쓴다.

```bash
minikube addons enable metrics-server
```

## 검증

```bash
kubectl get apiservice v1beta1.metrics.k8s.io
kubectl top nodes
kubectl top pods -A
```

## 제거

설치 방식에 맞는 명령 하나만 실행한다.

```bash
helm uninstall metrics-server -n kube-system
minikube addons disable metrics-server
```
