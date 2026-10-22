#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BENCH="$ROOT/benchmarks"
INPUT=${1:-"$ROOT/ad_data.csv"}
PROFILE=${2:-"$BENCH/ad_data.cpu.pprof"}
OUT="$BENCH/profile-results"
BIN="$BENCH/aggregator"
cd "$ROOT"

mkdir -p "$OUT"
go build -o "$BIN" "$ROOT/cmd/aggregator"
python3 "$BENCH/measure.py" --binary "$BIN" --input "$INPUT" --output "$OUT" \
	--cpuprofile "$PROFILE" --stdout "$BENCH/profile.stdout.log" \
	--stderr "$BENCH/profile.stderr.log" >"$BENCH/profile.runtime.metrics.log"
shasum -a 256 "$OUT/top10_ctr.csv" "$OUT/top10_cpa.csv" >"$BENCH/profile.output.sha256"
cmp "$BENCH/actual-results/top10_ctr.csv" "$OUT/top10_ctr.csv"
cmp "$BENCH/actual-results/top10_cpa.csv" "$OUT/top10_cpa.csv"

PPROF_KIND=
if go tool pprof -top "$BIN" "$PROFILE" >"$BENCH/profile.top.txt" 2>"$BENCH/profile.pprof.stderr.log"; then
	PPROF_KIND=go-tool
elif command -v pprof >/dev/null 2>&1; then
	pprof -top "$BIN" "$PROFILE" >"$BENCH/profile.top.txt" 2>"$BENCH/profile.pprof.stderr.log"
	PPROF_KIND=standalone
else
	printf '%s\n' "CPU profile collected; pprof analyzer unavailable. Install github.com/google/pprof separately to inspect it." >"$BENCH/profile.top.txt"
	printf '%s\n' "go tool pprof is not installed in this Go toolchain; standalone pprof is not on PATH." >"$BENCH/profile.pprof.stderr.log"
fi

if [ -n "$PPROF_KIND" ] && command -v dot >/dev/null 2>&1; then
	if [ "$PPROF_KIND" = "go-tool" ]; then
		go tool pprof -svg "$BIN" "$PROFILE" >"$BENCH/profile.svg" 2>"$BENCH/profile.svg.stderr.log"
	else
		pprof -svg "$BIN" "$PROFILE" >"$BENCH/profile.svg" 2>"$BENCH/profile.svg.stderr.log"
	fi
	printf '%s\n' "Graphviz dot available; generated profile.svg" >"$BENCH/profile.graphviz.log"
elif [ -z "$PPROF_KIND" ]; then
	printf '%s\n' "pprof analyzer unavailable; skipped SVG call graph" >"$BENCH/profile.graphviz.log"
else
	printf '%s\n' "Graphviz dot unavailable; skipped SVG call graph" >"$BENCH/profile.graphviz.log"
fi
