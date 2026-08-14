# Kubernetes Practice Monorepo Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a runnable Go and Kubernetes monorepo that teaches core concepts, production resilience, an ecommerce MSA, local clusters, CI, EKS portability, and on-premises architecture.

**Architecture:** One Go module contains a small reusable HTTP helper, a Hello API, and five ecommerce services. Kubernetes application resources use Kustomize bases plus minikube, kind, and EKS overlays; third-party platform products remain optional Helm add-ons. Documentation and labs follow the same deployable resources instead of copying independent projects.

**Tech Stack:** Go 1.26, Go standard library, Docker, Kubernetes 1.34+, kubectl/Kustomize, minikube, kind 0.32, Helm 4, GitHub Actions.

## Global Constraints

- Local default is Docker Desktop with minikube; kind is for multi-node and CI exercises.
- Application manifests use Kustomize. Helm is only for third-party add-ons.
- Go application code uses only the standard library.
- Containers run as non-root with immutable image tags, dropped Linux capabilities, disabled privilege escalation, and RuntimeDefault seccomp.
- Final ecommerce workloads use two replicas, probes, resource requests and limits, topology spread, PDBs, and graceful termination.
- Secret values are never committed. Documentation must state that base64 is not encryption.
- Service mesh is an optional Istio lab and is never required for the default path.
- EKS and RKE2 resources must not create infrastructure during verification.
- On-premises RKE2 content is documentation only.

---

## File Map

- `internal/web/web.go`: shared JSON, health, and graceful HTTP server helpers.
- `apps/hello-api/main.go`: single-service learning API.
- `apps/ecommerce/*/main.go`: gateway, catalog, inventory, orders, and payments services.
- `build/Dockerfile`: one parameterized multi-stage image build for all services.
- `deploy/base/*`: portable Kubernetes resources.
- `deploy/overlays/*`: environment-only differences.
- `deploy/addons/*`: optional Helm values and install instructions.
- `clusters/*`: kind and eksctl cluster declarations.
- `docs/concepts/*`: ordered Kubernetes concept reference.
- `docs/architecture/*`: ecommerce, reliability, and environment boundaries.
- `docs/runbooks/*`: symptom-first incident response.
- `labs/*/README.md`: executable learning exercises.
- `tests/*.sh`: manifest, documentation, smoke, and resilience checks.
- `.github/workflows/verify.yaml`: tests, builds, kind deployment, and smoke verification.

---

### Task 1: Repository Contract and Verification Shell

**Files:**
- Create: `.gitignore`
- Create: `go.mod`
- Create: `Makefile`
- Create: `build/Dockerfile`
- Create: `tests/manifests.sh`
- Create: `tests/docs.sh`

**Interfaces:**
- Consumes: no project code.
- Produces: `make test`, `make images`, `make manifests`, `make docs-check`, and `docker build -f build/Dockerfile --build-arg APP=apps/hello-api --build-arg BINARY=hello-api .`.

- [ ] **Step 1: Add failing repository checks**

`tests/manifests.sh` must render `deploy/overlays/minikube`, `deploy/overlays/kind`, and `deploy/overlays/eks` with `kubectl kustomize`. `tests/docs.sh` must assert that `README.md`, `docs/learning-path.md`, `docs/concepts/11-reliability-ha-dr.md`, and all sixteen lab README files exist.

- [ ] **Step 2: Run checks and confirm missing-file failures**

Run: `bash tests/manifests.sh; bash tests/docs.sh`

Expected: non-zero exits because deploy overlays, README, concepts, and labs do not exist.

- [ ] **Step 3: Add the root contracts**

Use module `k8s-practice` with `go 1.26`. The Dockerfile must build `CGO_ENABLED=0` from `golang:1.26.5-alpine`, then copy the binary into `scratch` with UID/GID `65532`. Make targets call `go test ./...` and build `k8s-practice/hello-api:dev`, `k8s-practice/gateway:dev`, `k8s-practice/catalog:dev`, `k8s-practice/inventory:dev`, `k8s-practice/orders:dev`, and `k8s-practice/payments:dev`.

- [ ] **Step 4: Verify the Makefile and Dockerfile syntax**

Run: `make -n test images manifests docs-check`

Expected: commands print without Make syntax errors. The repository checks remain red until later tasks create their inputs.

- [ ] **Step 5: Commit**

```bash
git add .gitignore go.mod Makefile build/Dockerfile tests/manifests.sh tests/docs.sh
git commit -m "build: add repository verification contract"
```

### Task 2: Shared HTTP Runtime and Hello API

**Files:**
- Create: `internal/web/web_test.go`
- Create: `internal/web/web.go`
- Create: `apps/hello-api/main_test.go`
- Create: `apps/hello-api/main.go`

**Interfaces:**
- Produces: `web.JSON(http.ResponseWriter, int, any)`, `web.Health(service string, next http.Handler) http.Handler`, and `web.Run(addr string, handler http.Handler) error`.
- Produces: Hello routes `GET /`, `GET /healthz`, `GET /readyz`, and `GET /config`.

- [ ] **Step 1: Write failing web helper tests**

Test that `JSON` sets `application/json`, that `Health` returns `{"status":"ok","service":"hello-api"}` on both probe paths, and that other routes reach the wrapped handler.

- [ ] **Step 2: Confirm helper tests fail**

Run: `go test ./internal/web`

Expected: compile failure because helper functions do not exist.

- [ ] **Step 3: Implement only the shared helper contract**

`Run` must use `http.Server` read-header, read, write, and idle timeouts; listen for `SIGINT` and `SIGTERM`; and call `Shutdown` with a ten-second context.

- [ ] **Step 4: Write failing Hello handler tests**

Construct `newHandler(getenv func(string) string) http.Handler`, request `/`, and assert JSON contains service `hello-api`, version `dev`, and pod `test-pod`. Request `/config` and assert only `MESSAGE` is returned, never arbitrary environment values.

- [ ] **Step 5: Implement the Hello handler and entrypoint**

Read `PORT` with default `8080`; wrap the mux in `web.Health`; use `APP_VERSION`, `HOSTNAME`, and `MESSAGE` for public response fields.

- [ ] **Step 6: Run unit tests**

Run: `go test ./internal/web ./apps/hello-api`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/web apps/hello-api
git commit -m "feat: add resilient hello api"
```

### Task 3: Foundation Kubernetes Resources and Labs 01-06

**Files:**
- Create: `deploy/base/hello/{namespace,configmap,deployment,service,kustomization}.yaml`
- Create: `deploy/overlays/minikube/kustomization.yaml`
- Create: `deploy/overlays/kind/kustomization.yaml`
- Create: `deploy/overlays/eks/kustomization.yaml`
- Create: `clusters/kind/config.yaml`
- Create: `labs/01-pod/README.md`
- Create: `labs/02-deployment-service/README.md`
- Create: `labs/03-config-secret/README.md`
- Create: `labs/04-probe-resource/README.md`
- Create: `labs/05-job-cronjob/README.md`
- Create: `labs/06-storage-statefulset/README.md`

**Interfaces:**
- Consumes: image `k8s-practice/hello-api:dev` and its four routes.
- Produces: namespace `practice`, Deployment/Service `hello-api`, and three renderable overlays.

- [ ] **Step 1: Run the manifest check and retain the missing-overlay failure**

Run: `bash tests/manifests.sh`

Expected: FAIL because overlays are absent.

- [ ] **Step 2: Implement the portable Hello base**

The namespace enforces Restricted Pod Security Admission. The Deployment uses two replicas, port 8080, all three probes, requests `20m/32Mi`, limits `200m/128Mi`, a non-root read-only security context, disabled service-account token mounting, and rolling update `maxUnavailable: 0`, `maxSurge: 1`.

- [ ] **Step 3: Implement the overlays and local kind cluster**

Minikube and kind reference the Hello base. EKS references the same base without credentials or AWS infrastructure. The kind config declares one control-plane and two workers.

- [ ] **Step 4: Write executable labs 01-06**

Each README includes goal, prerequisites, commands, observations, failure injection, recovery, and cleanup. Lab 03 creates a Secret using `kubectl create secret generic` and explicitly warns that base64 is not encryption. Lab 06 uses a temporary non-root BusyBox StatefulSet and PVC instead of adding a database dependency.

- [ ] **Step 5: Verify all overlays render**

Run: `bash tests/manifests.sh`

Expected: PASS and no rendered image uses `:latest`.

- [ ] **Step 6: Commit**

```bash
git add deploy/base/hello deploy/overlays clusters/kind labs/01-pod labs/02-deployment-service labs/03-config-secret labs/04-probe-resource labs/05-job-cronjob labs/06-storage-statefulset
git commit -m "feat: add foundation kubernetes labs"
```

### Task 4: Ecommerce Service Contracts

**Files:**
- Create: `apps/ecommerce/catalog/{main.go,main_test.go}`
- Create: `apps/ecommerce/inventory/{main.go,main_test.go}`
- Create: `apps/ecommerce/payments/{main.go,main_test.go}`
- Create: `apps/ecommerce/orders/{main.go,main_test.go}`
- Create: `apps/ecommerce/gateway/{main.go,main_test.go}`

**Interfaces:**
- Catalog: `GET /products` returns products `pencil` and `notebook`.
- Inventory: `GET /inventory/{productID}` returns deterministic availability.
- Payments: `POST /payments` accepts `order_id` and positive `amount`; returns `approved`.
- Orders: `POST /orders` accepts `product_id`, positive `quantity`, and positive `amount`; calls inventory then payments with a two-second client timeout.
- Gateway: proxies `GET /api/products` to catalog and `POST /api/orders` to orders.

- [ ] **Step 1: Write catalog, inventory, and payment handler tests**

Assert successful JSON responses, method-not-allowed responses, and bad-request responses for invalid identifiers or amounts.

- [ ] **Step 2: Confirm the three packages fail to compile**

Run: `go test ./apps/ecommerce/catalog ./apps/ecommerce/inventory ./apps/ecommerce/payments`

Expected: compile failures because `newHandler` does not exist.

- [ ] **Step 3: Implement the three leaf services**

Each package exposes `newHandler() http.Handler`, wraps it with `web.Health`, listens on port 8080, and keeps only deterministic in-memory sample data.

- [ ] **Step 4: Write orders orchestration tests**

Use two `httptest.Server` dependencies. Assert inventory is called before payments, unavailable inventory returns 409, dependency 5xx returns 502, malformed input returns 400, and a successful order returns an ID plus status `confirmed`.

- [ ] **Step 5: Implement orders with bounded dependency calls**

Construct `newHandler(inventoryURL, paymentsURL string, client *http.Client) http.Handler`. Do not retry non-idempotent payments. Convert dependency timeouts to HTTP 504 and other dependency errors to HTTP 502.

- [ ] **Step 6: Write and implement gateway proxy tests**

Construct `newHandler(catalogURL, ordersURL string) http.Handler` with `httputil.NewSingleHostReverseProxy`. Assert only the two documented paths route and an unknown path returns 404.

- [ ] **Step 7: Run the full Go suite**

Run: `go test ./...`

Expected: PASS without third-party modules.

- [ ] **Step 8: Commit**

```bash
git add apps/ecommerce
git commit -m "feat: add ecommerce service contracts"
```

### Task 5: Production-Shaped Ecommerce Manifests

**Files:**
- Create: `deploy/base/ecommerce/namespace.yaml`
- Create: `deploy/base/ecommerce/{gateway,catalog,inventory,orders,payments}.yaml`
- Create: `deploy/base/ecommerce/network-policies.yaml`
- Create: `deploy/base/ecommerce/kustomization.yaml`
- Modify: `deploy/overlays/{minikube,kind,eks}/kustomization.yaml`

**Interfaces:**
- Consumes: the five HTTP contracts from Task 4.
- Produces: namespace `shop`, five Deployments, five Services, five PDBs, and least-privilege network paths.

- [ ] **Step 1: Extend `tests/manifests.sh` with HA and security assertions**

Render each overlay, then assert ecommerce Deployment count is five, PDB count is five, no image uses `latest`, and rendered output includes `runAsNonRoot`, `RuntimeDefault`, `topologySpreadConstraints`, `maxUnavailable`, and `automountServiceAccountToken: false`.

- [ ] **Step 2: Run the check and confirm it fails**

Run: `bash tests/manifests.sh`

Expected: FAIL because ecommerce resources do not exist.

- [ ] **Step 3: Add the five workloads**

Each service file contains Service, Deployment, and PDB. Use two replicas, hostname topology spread with `ScheduleAnyway` so a single-node minikube remains usable, `minAvailable: 1`, probes, resources, rolling update, ten-second termination grace, and the global security context. Gateway references `catalog:8080` and `orders:8080`; orders references `inventory:8080` and `payments:8080`.

- [ ] **Step 4: Add network isolation**

Start with ingress and egress default deny. Permit DNS to kube-system, external traffic to gateway, gateway to catalog/orders, and orders to inventory/payments. Leaf services receive only from their declared caller.

- [ ] **Step 5: Add the base to all overlays and rerun checks**

Run: `bash tests/manifests.sh`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add deploy/base/ecommerce deploy/overlays tests/manifests.sh
git commit -m "feat: add highly available ecommerce manifests"
```

### Task 6: Core Concepts, Architecture, and Runbooks

**Files:**
- Create: `docs/learning-path.md`
- Create: `docs/concepts/01-mental-model.md`
- Create: `docs/concepts/02-cluster-architecture.md`
- Create: `docs/concepts/03-pod-lifecycle.md`
- Create: `docs/concepts/04-workload-controllers.md`
- Create: `docs/concepts/05-service-networking.md`
- Create: `docs/concepts/06-config-secret.md`
- Create: `docs/concepts/07-resources-scheduling.md`
- Create: `docs/concepts/08-storage-stateful.md`
- Create: `docs/concepts/09-security-multitenancy.md`
- Create: `docs/concepts/10-observability.md`
- Create: `docs/concepts/11-reliability-ha-dr.md`
- Create: `docs/concepts/12-packaging-gitops.md`
- Create: `docs/concepts/13-service-mesh.md`
- Create: `docs/concepts/14-maintenance-upgrades.md`
- Create: `docs/architecture/{ecommerce-msa,reliability-model,local-eks-comparison}.md`
- Create: `docs/runbooks/{crashloop-pending,network-dns,oom-eviction,rollout-rollback,node-drain,backup-restore,upgrade}.md`
- Create: `docs/onprem/rke2-production-reference.md`

**Interfaces:**
- Consumes: real paths and commands implemented in Tasks 1-5.
- Produces: ordered learning path and symptom-first operational reference.

- [ ] **Step 1: Run `tests/docs.sh` and retain the missing-document failure**

Run: `bash tests/docs.sh`

Expected: FAIL.

- [ ] **Step 2: Write concept documents with runnable examples**

Every document begins with learning outcomes and ends with a checklist plus links to the relevant lab. Reliability explicitly separates cluster, node, workload, data, dependency HA and DR; defines SLI, SLO, RTO, RPO; and explains that PDB protects only voluntary Eviction API disruptions.

- [ ] **Step 3: Write architecture and runbooks**

Runbooks use the fixed order: symptom, immediate safety check, diagnosis commands, recovery, verification, prevention. The on-prem reference covers RKE2 three-server HA, API load balancer, Cilium, MetalLB, Traefik Gateway, existing CSI or Longhorn, Harbor, monitoring, etcd snapshot off-site retention, and restore rehearsal without installation automation.

- [ ] **Step 4: Verify documentation structure**

Run: `bash tests/docs.sh`

Expected: only labs 07-16 remain missing.

- [ ] **Step 5: Commit**

```bash
git add docs tests/docs.sh
git commit -m "docs: add kubernetes concepts and runbooks"
```

### Task 7: Platform and Advanced Labs 07-16

**Files:**
- Create: `deploy/addons/gateway/{README.md,values.yaml}`
- Create: `deploy/addons/metrics/README.md`
- Create: `deploy/addons/observability/{README.md,values.yaml,rules.yaml}`
- Create: `deploy/addons/gitops/{README.md,values.yaml}`
- Create: `deploy/addons/service-mesh/{README.md,values.yaml}`
- Create: `labs/07-ingress-gateway/README.md`
- Create: `labs/08-autoscaling-scheduling/README.md`
- Create: `labs/09-rbac-network-policy/README.md`
- Create: `labs/10-observability/README.md`
- Create: `labs/11-rollout-rollback/README.md`
- Create: `labs/12-failure-recovery/README.md`
- Create: `labs/13-ecommerce-msa/README.md`
- Create: `labs/14-service-mesh/README.md`
- Create: `labs/15-gitops/README.md`
- Create: `labs/16-eks-migration/README.md`

**Interfaces:**
- Consumes: Helm 4 and resources from Tasks 3 and 5.
- Produces: optional Traefik Gateway API, Metrics Server, kube-prometheus-stack/Alertmanager, Argo CD, and Istio workflows.

- [ ] **Step 1: Add failing addon checks to `tests/docs.sh`**

Require each add-on README and values file listed above, and require all sixteen lab README files.

- [ ] **Step 2: Write add-on contracts**

Each add-on README includes repository setup, install, verify, and uninstall commands. Observability uses kube-prometheus-stack and alert rules linked to repository runbooks. Istio is namespace-scoped to `shop` and remains optional.

- [ ] **Step 3: Write labs 07-16**

Each lab uses the common format from Tasks 3 and 6. Lab 12 covers Pod deletion, drain, OOM, DNS, NetworkPolicy, dependency timeout, bad image rollout, PDB blocking, and capacity shortage with native tools. Lab 16 renders EKS resources but clearly marks AWS commands as cost-incurring and user-initiated.

- [ ] **Step 4: Run documentation checks**

Run: `bash tests/docs.sh`

Expected: PASS with no empty files and no unresolved internal Markdown paths.

- [ ] **Step 5: Commit**

```bash
git add deploy/addons labs tests/docs.sh
git commit -m "docs: add platform and advanced labs"
```

### Task 8: Local Workflow, EKS Overlay, Smoke Tests, and CI

**Files:**
- Create: `clusters/eks/cluster.yaml`
- Create: `deploy/overlays/eks/ingress.yaml`
- Create: `tests/smoke.sh`
- Create: `tests/resilience.sh`
- Create: `.github/workflows/verify.yaml`
- Create: `README.md`
- Modify: `Makefile`
- Modify: `deploy/overlays/eks/kustomization.yaml`

**Interfaces:**
- Produces: `make minikube-up`, `make kind-up`, `make load-images`, `make deploy`, `make smoke`, `make resilience-kind`, and `make clean-kind`.
- Produces: an eksctl configuration for `ap-northeast-2` with three managed worker nodes and no automatic invocation.

- [ ] **Step 1: Write the failing smoke script**

Port-forward gateway service to a random local port, assert `/api/products` contains `pencil`, post a valid order, and assert status `confirmed`. The script must clean up the port-forward process on exit.

- [ ] **Step 2: Add local workflow targets and CI**

CI installs Go 1.26.x and kind 0.32.0, runs unit and documentation checks, builds six images, creates the three-node kind cluster, loads images, applies the kind overlay, waits for rollouts, and runs the smoke script. Cleanup runs even after failure.

- [ ] **Step 3: Add resilience verification**

Delete one gateway Pod and wait for availability, drain one worker with the Eviction API, assert a ready gateway remains due to replicas/PDB, uncordon the worker, and wait for full recovery. Abrupt node stop remains a manual lab because node eviction detection is intentionally slow.

- [ ] **Step 4: Add EKS-specific boundaries**

The EKS overlay adds an ALB Ingress for gateway and documents the required AWS Load Balancer Controller, ECR image rewrite, EBS CSI boundary, and Pod Identity boundary. Verification only renders YAML.

- [ ] **Step 5: Write the root README**

Include prerequisites, quick start, tool roles, repository map, learning tracks, common Make commands, local versus EKS versus on-prem boundaries, cost warning, troubleshooting entry points, and GitHub remote commands without assuming a remote URL.

- [ ] **Step 6: Run local non-cluster verification**

Run: `go test ./... && make manifests docs-check && git diff --check`

Expected: PASS.

- [ ] **Step 7: Run kind verification when Docker is available**

Run: `make kind-up images load-images deploy smoke resilience-kind`

Expected: all Deployments available, smoke PASS, Pod deletion recovery PASS, drain recovery PASS.

- [ ] **Step 8: Commit**

```bash
git add README.md Makefile clusters/eks deploy/overlays/eks tests .github/workflows/verify.yaml
git commit -m "ci: verify local and portable kubernetes workflows"
```

### Task 9: Final Consistency and Security Review

**Files:**
- Modify: only files proven inconsistent by verification.

**Interfaces:**
- Consumes: every command and path documented in Tasks 1-8.
- Produces: clean repository with a runnable learning path.

- [ ] **Step 1: Run complete static verification**

Run: `go test ./...`, `make manifests`, `make docs-check`, `git diff --check`, and `rg -n 'image:.*:latest|privileged: true|runAsNonRoot: false' deploy`.

Expected: tests/checks PASS and the security search returns no matches.

- [ ] **Step 2: Verify all documented relative paths and Make targets**

Run: `bash tests/docs.sh && make -n minikube-up kind-up images load-images deploy smoke resilience-kind clean-kind`

Expected: PASS without missing targets.

- [ ] **Step 3: Run Docker/kind verification or record the exact environment blocker**

If the Docker daemon and kind are available, run Task 8 Step 7. If unavailable, keep the static and container-build evidence and report the missing executable or daemon state exactly; do not claim cluster verification passed.

- [ ] **Step 4: Commit verification-only corrections**

```bash
git add -A
git commit -m "chore: finish kubernetes practice verification"
```
