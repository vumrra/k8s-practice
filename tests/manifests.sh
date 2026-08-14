#!/bin/sh
set -eu

for overlay in minikube kind eks; do
  kubectl kustomize "deploy/overlays/$overlay" >/dev/null
done

