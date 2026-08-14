# Lab 05 — Job과 CronJob

## 목표

계속 실행되는 Deployment와 완료를 목표로 하는 Job·주기 실행 CronJob을 구분한다.

## 준비

`practice` namespace와 인터넷에서 `busybox:1.36.1` 이미지를 가져올 수 있는 클러스터가 필요하다.

## Job

```yaml
# /tmp/practice-job.yaml
apiVersion: batch/v1
kind: Job
metadata: {name: once, namespace: practice}
spec:
  template:
    spec:
      restartPolicy: Never
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        seccompProfile: {type: RuntimeDefault}
      containers:
        - name: task
          image: busybox:1.36.1
          command: ["sh", "-c", "date; echo completed"]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities: {drop: ["ALL"]}
```

```bash
kubectl apply -f /tmp/practice-job.yaml
kubectl wait -n practice --for=condition=complete job/once --timeout=60s
kubectl logs -n practice job/once
```

CronJob은 같은 Job template에 `schedule: "*/5 * * * *"`, `concurrencyPolicy: Forbid`, `successfulJobsHistoryLimit`, `failedJobsHistoryLimit`을 추가한다. 중복 실행이 위험한 작업은 애플리케이션도 idempotent해야 한다.

## 실패와 정리

존재하지 않는 명령으로 바꾸면 Job이 backoffLimit까지 재시도한다. 원인을 고친 새 Job을 만들며 완료된 Pod를 재사용하지 않는다.

```bash
kubectl delete -f /tmp/practice-job.yaml --ignore-not-found
```
