# 03. Pod와 수명주기

## 학습 목표

- Pod의 정체성·수명주기와 container restart를 구분한다.
- init container, sidecar, probe, graceful termination을 적용할 시점을 안다.

Pod는 하나 이상의 밀접한 컨테이너가 network namespace와 volume을 공유하는 최소 배포 단위다. Pod IP와 UID는 재생성 시 바뀌므로 영구 식별자로 사용하지 않는다.

## 상태와 재시작

- Pending: 스케줄링, 이미지, volume 준비 전
- Running: 하나 이상의 container 실행 중
- Succeeded/Failed: 종료된 batch Pod
- Unknown: Node와 통신 불가
- CrashLoopBackOff: 상태가 아니라 반복 실패에 대한 재시작 backoff 표시

`restartPolicy`는 같은 Pod 안의 container 재시작 규칙이다. Node 장애처럼 Pod 자체가 사라지면 Deployment·StatefulSet 같은 상위 controller가 새 Pod를 만든다.

## 특수 container

- init container: 본 container 전에 순서대로 완료되어야 하는 준비 작업
- sidecar: 같은 Pod에서 주 container의 수명주기를 보조하는 장기 실행 container
- ephemeral container: 실행 중 Pod 디버깅용이며 일반 workload 정의가 아님

한 Pod에 넣으면 독립 배포와 확장이 불가능해진다. 항상 함께 시작·종료하고 localhost/volume을 공유해야 할 때만 묶는다.

## Probe

- startup: 느린 초기화가 끝났는지 확인
- liveness: 멈춘 process를 재시작할지 결정
- readiness: 현재 Service 트래픽을 받을 수 있는지 결정

외부 DB 장애 때문에 liveness를 실패시키면 재시작 폭풍이 생길 수 있다. dependency 상태는 보통 readiness와 애플리케이션 timeout으로 처리한다.

## 종료

삭제 시 kubelet은 종료 유예를 시작하고 endpoint에서 제외한 뒤 SIGTERM을 전달한다. 애플리케이션은 새 요청을 거부하고 진행 중 요청을 마친 후 종료해야 한다. 유예가 끝나면 SIGKILL된다.

```bash
kubectl get pod -n practice -w
kubectl describe pod -n practice POD_NAME
kubectl logs -n practice POD_NAME --previous
kubectl exec -n practice POD_NAME -- COMMAND
```

## 체크리스트

- [ ] Pod IP와 로컬 filesystem을 영구 상태로 사용하지 않는다.
- [ ] readiness와 liveness 목적을 바꾸어 쓰지 않는다.
- [ ] SIGTERM과 `terminationGracePeriodSeconds`를 함께 설계한다.

실습: [Lab 01](../../labs/01-pod/README.md), [Lab 04](../../labs/04-probe-resource/README.md)
