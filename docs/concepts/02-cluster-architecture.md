# 02. 클러스터 아키텍처

## 학습 목표

- control plane과 worker node의 책임을 구분한다.
- Pod 생성 요청이 실제 컨테이너가 되기까지의 흐름을 설명한다.

## 구성요소

```text
kubectl
  -> kube-apiserver: 인증·인가·admission과 API 진입점
      -> etcd: 클러스터 상태 저장
      -> scheduler: 미배치 Pod에 Node 선택
      -> controller-manager: desired/current state 조정

worker node
  -> kubelet: PodSpec대로 컨테이너 실행·상태 보고
  -> containerd 같은 runtime: 이미지와 컨테이너 관리
  -> CNI: Pod 네트워크
  -> kube-proxy 또는 대체 dataplane: Service 트래픽
  -> CSI node plugin: volume 연결
```

API Server가 컨테이너를 직접 실행하지 않는다. scheduler는 Node를 선택하고 kubelet이 runtime을 통해 실행한다. controller와 kubelet은 장애 후에도 반복해서 상태를 맞춘다.

## 로컬과 운영의 차이

minikube 단일 노드는 학습에 적합하지만 node/control plane 장애를 견디지 못한다. kind 다중 노드는 스케줄링과 drain을 실습할 수 있지만 Docker host 하나가 공통 장애 지점이다. 운영 HA는 별도 서버·가용 영역, 다중 control plane, API load balancer와 etcd quorum을 필요로 한다.

## 확인

```bash
kubectl cluster-info
kubectl get componentstatuses  # 일부 배포판에서는 제한됨
kubectl get node -o wide
kubectl get pod -n kube-system -o wide
kubectl describe node
```

## 체크리스트

- [ ] etcd를 정기 백업해야 하는 이유를 안다.
- [ ] worker 장애와 control plane 장애의 영향을 구분한다.
- [ ] 로컬 다중 노드가 실제 인프라 HA는 아님을 안다.

실습: [Lab 12](../../labs/12-failure-recovery/README.md) · [공식 컴포넌트 문서](https://kubernetes.io/docs/concepts/overview/components/)
