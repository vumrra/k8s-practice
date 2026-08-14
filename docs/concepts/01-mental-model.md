# 01. Kubernetes 사고방식과 Object

## 학습 목표

- 명령형 실행과 선언형 desired state의 차이를 설명한다.
- `spec`, `status`, label, selector, annotation, owner reference, namespace를 구분한다.

## 핵심 모델

Kubernetes Object는 “무엇을 실행하라”는 일회성 명령이 아니라 원하는 상태의 기록이다.

```text
사용자 YAML(spec) -> API Server -> 저장된 desired state
                              -> controller가 현재 상태 관찰
                              -> 차이를 줄이는 reconcile 반복
                              -> status 갱신
```

Pod를 직접 지우면 끝이지만 Deployment가 소유한 Pod를 지우면 ReplicaSet controller가 새 Pod를 만든다. 이것이 자기 복구의 기본이며, 애플리케이션 오류를 자동으로 고친다는 뜻은 아니다.

## Object의 공통 필드

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-api
  namespace: practice
  labels:
    app.kubernetes.io/name: hello-api
spec: {}
status: {}
```

- name/UID: object 식별
- label: selector와 검색에 사용하는 식별 속성
- annotation: controller나 도구에 전달하는 비식별 메타데이터
- ownerReference: 상위 controller와 garbage collection 관계
- namespace: 이름과 권한·quota 범위. 보안 경계로만 믿어서는 안 된다.

## 자주 쓰는 관찰 명령

```bash
kubectl api-resources
kubectl explain deployment.spec.strategy
kubectl get deployment -n practice -o yaml
kubectl get pod -A -l app.kubernetes.io/name=hello-api
kubectl diff -k deploy/overlays/minikube
kubectl apply -k deploy/overlays/minikube
```

`apply`는 같은 desired state를 반복 적용할 수 있어야 한다. 긴급 수정을 `kubectl edit`로 했다면 Git의 선언과 달라진 drift를 다시 정리한다.

## 체크리스트

- [ ] `spec`은 사용자 의도, `status`는 controller가 관찰한 현재 상태라고 설명한다.
- [ ] label과 annotation의 용도를 구분한다.
- [ ] 기본 namespace 대신 workload 전용 namespace를 사용한다.

실습: [Lab 01](../../labs/01-pod/README.md), [Lab 02](../../labs/02-deployment-service/README.md) · [공식 Object 문서](https://kubernetes.io/docs/concepts/overview/working-with-objects/)
