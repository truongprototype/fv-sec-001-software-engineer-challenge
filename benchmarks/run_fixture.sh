#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BENCH="$ROOT/benchmarks"
OUT="$BENCH/fixture-results"
BIN="$BENCH/aggregator"
cd "$ROOT"

mkdir -p "$OUT"
go build -o "$BIN" "$ROOT/cmd/aggregator"
"$BIN" --input "$BENCH/fixture.csv" --output "$OUT" >"$BENCH/fixture.stdout.log" 2>"$BENCH/fixture.stderr.log"
diff -u "$BENCH/expected_top10_ctr.csv" "$OUT/top10_ctr.csv"
diff -u "$BENCH/expected_top10_cpa.csv" "$OUT/top10_cpa.csv"
echo "fixture verification: PASS" >"$BENCH/fixture.verification.log"
