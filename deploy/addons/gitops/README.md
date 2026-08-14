# Argo CD

Git의 선언 상태와 클러스터 상태를 지속적으로 비교·동기화한다. 이 설정은 로컬 단일 복제본 실습용이며 운영 HA 구성은 Argo CD 지원 정책에 맞춰 별도로 잡는다.

## 설치

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
helm upgrade --install argocd argo/argo-cd \
  --namespace argocd --create-namespace \
  --values deploy/addons/gitops/values.yaml --wait
```

## 검증

```bash
kubectl get pods -n argocd
kubectl port-forward -n argocd svc/argocd-server 8080:443
```

초기 비밀번호를 출력하는 명령은 화면 공유·셸 기록에 노출될 수 있다.

```bash
kubectl get secret -n argocd argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | base64 --decode
```

Application 생성은 [실습 15](../../../labs/15-gitops/README.md)에서 저장소 URL을 정한 뒤 진행한다. 비밀키나 저장소 자격 증명을 Git에 평문으로 넣지 않는다.

## 제거

```bash
helm uninstall argocd -n argocd
kubectl delete namespace argocd
```
