#!/bin/sh
set -eu

render_dir=$(mktemp -d)
trap 'rm -rf "$render_dir"' EXIT

for overlay in minikube kind eks; do
  rendered="$render_dir/$overlay.yaml"
  kubectl kustomize "deploy/overlays/$overlay" >"$rendered"

  deployments=$(grep -c '^kind: Deployment$' "$rendered")
  pdbs=$(grep -c '^kind: PodDisruptionBudget$' "$rendered" || true)
  test "$deployments" -eq 6 || { echo "$overlay: expected 6 Deployments, found $deployments" >&2; exit 1; }
  test "$pdbs" -eq 5 || { echo "$overlay: expected 5 PDBs, found $pdbs" >&2; exit 1; }

  grep -q 'runAsNonRoot: true' "$rendered"
  grep -q 'type: RuntimeDefault' "$rendered"
  grep -q 'topologySpreadConstraints:' "$rendered"
  grep -q 'maxUnavailable: 0' "$rendered"
  grep -q 'automountServiceAccountToken: false' "$rendered"
  if grep -Eq 'image: .*:latest([[:space:]]|$)' "$rendered"; then
    echo "$overlay: latest image tag is forbidden" >&2
    exit 1
  fi
done
