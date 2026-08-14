#!/bin/sh
set -eu

for file in \
  README.md \
  docs/learning-path.md \
  docs/concepts/11-reliability-ha-dr.md \
  deploy/addons/gateway/README.md \
  deploy/addons/gateway/values.yaml \
  deploy/addons/metrics/README.md \
  deploy/addons/observability/README.md \
  deploy/addons/observability/values.yaml \
  deploy/addons/observability/rules.yaml \
  deploy/addons/gitops/README.md \
  deploy/addons/gitops/values.yaml \
  deploy/addons/service-mesh/README.md \
  deploy/addons/service-mesh/values.yaml; do
  test -s "$file" || { echo "missing: $file" >&2; exit 1; }
done

lab_count=$(find labs -mindepth 2 -maxdepth 2 -name README.md 2>/dev/null | wc -l | tr -d ' ')
test "$lab_count" -eq 16 || { echo "expected 16 labs, found $lab_count" >&2; exit 1; }

if grep -R -n -E -- '-l app=' tests/smoke.sh tests/resilience.sh labs deploy/addons; then
  echo "use the app.kubernetes.io/name label selector" >&2
  exit 1
fi

grep -q 'istio-cni istio/cni' deploy/addons/service-mesh/README.md || {
  echo "service mesh lab must use Istio CNI with restricted Pod Security" >&2
  exit 1
}

find README.md docs deploy labs -type f -name '*.md' -exec sh -c '
  for file do
    base=$(dirname "$file")
    sed -n "s/.*](\([^)]*\)).*/\1/p" "$file" |
      while IFS= read -r link; do
        case "$link" in
          ""|\#*|http://*|https://*|mailto:*) continue ;;
        esac
        target=${link%%#*}
        test -e "$base/$target" || {
          echo "broken link: $file -> $link" >&2
          exit 1
        }
      done
  done
' sh {} +
