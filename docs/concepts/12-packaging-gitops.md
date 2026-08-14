# 12. Kustomize, Helm과 GitOps

## 학습 목표

- raw YAML, Kustomize와 Helm의 책임을 구분한다.
- GitOps reconciliation과 배포 pipeline의 차이를 설명한다.

## 도구 선택

| 방식 | 적합한 경우 |
|---|---|
| raw YAML | object 몇 개를 처음 학습 |
| Kustomize | 같은 앱의 환경별 patch/image 차이 |
| Helm | 배포 가능한 제품 패키지와 외부 add-on |

이 저장소는 직접 만든 앱은 Kustomize, Traefik·Prometheus·Argo CD·Istio는 공식 Helm chart를 쓴다. 앱을 Helm과 Kustomize 양쪽에서 중복 관리하지 않는다.

```text
deploy/base -> minikube overlay
            -> kind overlay
            -> EKS overlay
```

```bash
kubectl kustomize deploy/overlays/kind
kubectl diff -k deploy/overlays/kind
kubectl apply -k deploy/overlays/kind
helm upgrade --install RELEASE CHART -n NAMESPACE --create-namespace -f values.yaml
helm rollback RELEASE REVISION -n NAMESPACE
```

## GitOps

```text
Git desired state -> Argo CD가 관찰 -> cluster current state와 비교
                 -> sync/self-heal/prune 정책에 따라 reconcile
```

CI가 cluster credential로 직접 push하는 방식과 달리 GitOps controller가 cluster 안에서 pull/reconcile한다. 긴급 수동 수정은 Git과 drift가 생기며 self-heal에 의해 되돌아갈 수 있다.

Secret plaintext를 GitOps 저장소에 넣지 않는다. external secret manager, SOPS/Sealed Secrets 같은 별도 선택이 필요하지만 이 저장소는 제품 하나를 기본 강제하지 않는다.

## 체크리스트

- [ ] 렌더링 결과를 review하고 적용한다.
- [ ] chart와 image version을 운영에서 고정한다.
- [ ] prune 전에 삭제 영향과 CRD 순서를 확인한다.
- [ ] rollback과 DB migration 호환성을 함께 설계한다.

실습: [Lab 15](../../labs/15-gitops/README.md)
