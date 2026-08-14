#!/bin/sh
set -eu

port=$((18080 + ($$ % 1000)))
log_file=$(mktemp)

cleanup() {
  if test -n "${port_forward_pid:-}"; then
    kill "$port_forward_pid" 2>/dev/null || true
    wait "$port_forward_pid" 2>/dev/null || true
  fi
  rm -f "$log_file"
}
trap cleanup EXIT INT TERM

kubectl port-forward -n shop service/gateway "$port:8080" >"$log_file" 2>&1 &
port_forward_pid=$!

attempt=0
until curl --fail --silent "http://127.0.0.1:$port/readyz" >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if test "$attempt" -ge 30 || ! kill -0 "$port_forward_pid" 2>/dev/null; then
    cat "$log_file" >&2
    echo "gateway port-forward did not become ready" >&2
    exit 1
  fi
  sleep 1
done

products=$(curl --fail --silent --show-error "http://127.0.0.1:$port/api/products")
echo "$products" | grep -q '"id":"pencil"' || {
  echo "catalog response does not contain pencil: $products" >&2
  exit 1
}

order=$(curl --fail --silent --show-error \
  -X POST "http://127.0.0.1:$port/api/orders" \
  -H 'Content-Type: application/json' \
  -d '{"product_id":"pencil","quantity":2,"amount":3}')
echo "$order" | grep -q '"status":"confirmed"' || {
  echo "order was not confirmed: $order" >&2
  exit 1
}

echo "smoke: PASS"
