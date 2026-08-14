# 06. ConfigMap과 Secret

## 학습 목표

- 일반 설정과 민감 정보를 분리한다.
- 환경 변수와 volume mount 갱신 방식의 차이를 안다.

ConfigMap은 공개 가능한 설정, Secret은 password·token·key 같은 민감 정보의 Kubernetes 전달 형식이다. 둘 다 애플리케이션 이미지와 환경별 값을 분리한다.

```yaml
env:
  - name: MESSAGE
    valueFrom:
      configMapKeyRef: {name: hello-api, key: MESSAGE}
```

- env: process 시작 시 읽으므로 object가 바뀌어도 재시작 전에는 변하지 않는다.
- volume: kubelet이 파일을 갱신하지만 앱이 다시 읽어야 한다. `subPath` mount는 자동 갱신되지 않는다.
- immutable object: 실수로 변경되는 것을 막고 watch 부하를 줄일 수 있지만 이름 교체가 필요하다.

## Secret 주의점

base64는 암호화가 아니다. Secret은 기본 설정에서 etcd에 암호화되지 않을 수 있다. Git에는 실제 값을 넣지 말고 다음을 함께 설계한다.

- etcd encryption at rest 또는 관리형 서비스 설정 확인
- `get/list/watch`를 최소화한 RBAC
- Secret이 필요 없는 container에는 mount하지 않기
- 로그·오류 응답에 값 노출 금지
- 짧은 수명, rotation과 폐기 절차
- 필요하면 외부 secret manager와 CSI/동기화 도구 사용

```bash
kubectl create secret generic -n practice demo --from-literal=TOKEN=practice-only
kubectl auth can-i get secrets -n practice --as=system:serviceaccount:practice:default
```

## 체크리스트

- [ ] ConfigMap에 credential을 넣지 않는다.
- [ ] Secret YAML을 Git에 저장하지 않는다.
- [ ] workload 생성 권한이 Secret 우회 열람으로 이어질 수 있음을 안다.

실습: [Lab 03](../../labs/03-config-secret/README.md) · [공식 Secret 보안 문서](https://kubernetes.io/docs/concepts/security/secrets-good-practices/)
