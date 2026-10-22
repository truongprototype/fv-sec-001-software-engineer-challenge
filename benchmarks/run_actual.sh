#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BENCH="$ROOT/benchmarks"
INPUT=${1:-"$ROOT/ad_data.csv"}
OUT="$BENCH/actual-results"
BIN="$BENCH/aggregator"
cd "$ROOT"

mkdir -p "$OUT"
go version >"$BENCH/environment.log"
python3 --version >>"$BENCH/environment.log" 2>&1
go env GOOS GOARCH GOMOD >>"$BENCH/environment.log"
uname -a >>"$BENCH/environment.log"
python3 "$BENCH/environment.py" >>"$BENCH/environment.log"
stat -f 'input_bytes=%z' "$INPUT" >>"$BENCH/environment.log"
shasum -a 256 "$INPUT" >>"$BENCH/environment.log"
go build -o "$BIN" "$ROOT/cmd/aggregator"
go run "$BENCH/count.go" --input "$INPUT" >"$BENCH/count.stdout.log" 2>"$BENCH/count.stderr.log"
python3 "$BENCH/measure.py" --binary "$BIN" --input "$INPUT" --output "$OUT" --stdout "$BENCH/actual.stdout.log" --stderr "$BENCH/actual.stderr.log" >"$BENCH/runtime.metrics.log"
shasum -a 256 "$OUT/top10_ctr.csv" "$OUT/top10_cpa.csv" >"$BENCH/output.sha256"
