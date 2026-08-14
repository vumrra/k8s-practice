# 15. Argo CD GitOps

## 목표

GitHub에 올린 이 저장소의 kind overlay를 Argo CD가 비교하도록 연결하고 drift를 확인한다.

## 준비

[Argo CD 애드온](../../deploy/addons/gitops/README.md)을 설치하고 저장소를 GitHub에 push한다. 비공개 저장소 인증은 UI나 Secret 관리 체계로 별도 등록한다.

## 실행

자신의 저장소 URL로 바꾼다.

```bash
export REPO_URL=https://github.com/OWNER/k8s-practice.git
kubectl apply -f - <<YAML
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: k8s-practice
  namespace: argocd
spec:
  project: default
  source:
    repoURL: ${REPO_URL}
    targetRevision: main
    path: deploy/overlays/kind
  destination:
    server: https://kubernetes.default.svc
    namespace: shop
  syncPolicy:
    syncOptions:
      - CreateNamespace=true
YAML
kubectl get application -n argocd k8s-practice -w
```

## 관찰

```bash
kubectl describe application -n argocd k8s-practice
kubectl scale deployment -n shop gateway --replicas=1
kubectl get application -n argocd k8s-practice
```

수동 변경은 `OutOfSync`를 만든다. 먼저 차이를 읽고 동기화한다. 자동 prune/self-heal은 영향 범위를 이해한 뒤 켠다.

## 장애와 복구

잘못된 Git revision이나 경로를 지정해 `ComparisonError`를 관찰하고 원래 값으로 복구한다. GitOps 컨트롤러 장애 시 이미 실행 중인 워크로드는 계속 동작하지만 새 배포는 멈춘다.

## 정리

```bash
kubectl delete application -n argocd k8s-practice
kubectl apply -k deploy/overlays/kind
```
