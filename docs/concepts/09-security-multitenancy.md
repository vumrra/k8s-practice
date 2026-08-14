# 09. 보안과 Multi-tenancy

## 학습 목표

- authentication, authorization, admission의 순서를 안다.
- workload와 namespace에 최소 권한·격리를 적용한다.

```text
API 요청 -> Authentication -> Authorization(RBAC) -> Admission -> etcd
```

- ServiceAccount: Pod의 Kubernetes API identity
- Role/ClusterRole: 허용 동작 집합
- RoleBinding/ClusterRoleBinding: identity와 권한 연결
- Admission: 저장 전에 정책 검증·변경
- audit log: 누가 언제 어떤 API 작업을 했는지 기록

## Workload 기본값

```yaml
automountServiceAccountToken: false
securityContext:
  runAsNonRoot: true
  seccompProfile: {type: RuntimeDefault}
containers:
  - securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities: {drop: ["ALL"]}
```

namespace에는 Pod Security Admission `restricted`를 적용한다. privileged Pod, hostPath, hostNetwork와 broad RBAC는 node·cluster 탈취 경로가 될 수 있다.

## 공유 클러스터

- 신뢰 수준·team별 namespace
- least privilege RBAC
- default deny NetworkPolicy
- ResourceQuota·LimitRange
- 외부 identity와 짧은 credential
- 이미지 digest/tag 고정, 취약점·서명 검사와 registry 정책

namespace는 강한 tenant 경계가 아니다. workload 생성 권한은 같은 namespace의 Secret과 ServiceAccount를 간접 사용할 수 있으므로 민감도가 다른 workload를 분리한다.

```bash
kubectl auth can-i --list -n shop
kubectl auth can-i create pods -n shop --as=USER
kubectl get namespace shop --show-labels
kubectl get networkpolicy -n shop
```

## 체크리스트

- [ ] API를 쓰지 않는 Pod에서 token mount를 끈다.
- [ ] `cluster-admin`과 wildcard RBAC를 애플리케이션에 주지 않는다.
- [ ] Secret base64, image `latest`, root container를 운영 기본값으로 쓰지 않는다.

실습: [Lab 09](../../labs/09-rbac-network-policy/README.md) · [공식 보안 개요](https://kubernetes.io/docs/concepts/security/)
