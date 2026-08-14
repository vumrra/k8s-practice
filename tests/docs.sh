#!/bin/sh
set -eu

for file in \
  README.md \
  docs/learning-path.md \
  docs/concepts/11-reliability-ha-dr.md; do
  test -s "$file" || { echo "missing: $file" >&2; exit 1; }
done

lab_count=$(find labs -mindepth 2 -maxdepth 2 -name README.md 2>/dev/null | wc -l | tr -d ' ')
test "$lab_count" -eq 16 || { echo "expected 16 labs, found $lab_count" >&2; exit 1; }

