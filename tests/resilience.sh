#!/bin/sh
set -eu

expected_context=kind-k8s-practice
node=k8s-practice-worker

context=$(kubectl config current-context)
test "$context" = "$expected_context" || {
  echo "refusing resilience test on context '$context'; expected '$expected_context'" >&2
  exit 1
}
kubectl get node "$node" >/dev/null

uncordon() {
  kubectl uncordon "$node" >/dev/null 2>&1 || true
}
trap uncordon EXIT INT TERM

pod=$(kubectl get pod -n shop -l app.kubernetes.io/name=gateway -o jsonpath='{.items[0].metadata.name}')
kubectl delete pod -n shop "$pod" --wait=false
kubectl rollout status deployment/gateway -n shop --timeout=180s

kubectl drain "$node" --ignore-daemonsets --delete-emptydir-data --timeout=180s
ready=$(kubectl get deployment/gateway -n shop -o jsonpath='{.status.readyReplicas}')
test "${ready:-0}" -ge 1 || {
  echo "gateway lost all ready replicas during drain" >&2
  exit 1
}

uncordon
kubectl rollout status deployment/gateway -n shop --timeout=180s
echo "resilience: PASS"
