# Benchmark record

## Correctness fixture

`fixture.csv` is a six-row hand-written fixture. The expected Top CTR and Top CPA files are checked in as `expected_top10_ctr.csv` and `expected_top10_cpa.csv`. The fixture covers aggregation across multiple rows, zero impressions, zero conversions, deterministic ordering, CSV output, and spend precision beyond cents.

Reproduce the fixture verification from the repository root:

```sh
sh benchmarks/run_fixture.sh
```

The script builds the binary, runs the fixture, and diffs both outputs against the expected CSV files. The verified run is recorded in `fixture.verification.log`; application stdout/stderr are in `fixture.stdout.log` and `fixture.stderr.log`.

## Actual input benchmark

Measured on 2026-10-23. `ad_data.csv` was already present; no data generation or decompression occurred during the application run.

| Property | Value |
|---|---:|
| Input size | 1,043,304,870 bytes (about 995 MiB) |
| SHA-256 | `4e9c87cc0b408f043917f5f7ee0101ba09899ad79f154358e9a2da374cfd60eb` |
| Data rows (header excluded) | 26,843,544 |
| Distinct campaigns | 50 |
| OS/kernel | macOS Darwin 27.0.0, arm64 |
| CPU info available to process | ARM, 8 logical CPUs; exact processor model was not exposed |
| Installed memory | 17,179,869,184 bytes (16 GiB) |
| Go | go1.27.3 darwin/arm64 |
| Python used by measurement helper | Python 3.9.6 |

Two successful application runs used the same input and emitted byte-identical outputs:

| Run | Wall time | User CPU | System CPU | Peak RSS |
|---|---:|---:|---:|---:|
| 1 | 21.302 s | not captured | not captured | 12,959,744 bytes |
| 2 | 21.342 s | 21.327 s | 0.674 s | 13,287,424 bytes |

Peak RSS comes from `resource.getrusage(RUSAGE_CHILDREN).ru_maxrss` in `measure.py`; on macOS this is bytes. Wall time covers only the aggregator subprocess, including opening, parsing, aggregating, ranking, and writing both CSVs. It excludes binary build, input hashing, the separate row/campaign counting pass, fixture creation, and extraction. The count helper and environment/hash collection are separately logged.

The measured app used about 12.7 MiB peak RSS while consuming a 995 MiB input, consistent with streaming rows and retaining state by campaign. User CPU time was nearly equal to wall time on run 2, suggesting a CPU-bound pass. With only 50 campaigns, sorting is unlikely to dominate; CSV/date/decimal parsing and exact rational arithmetic are plausible cost centers. This is an inference, not a profile result. No application optimization was made. Profile before selecting a change, then rerun all tests and compare output hashes against this baseline.

Reproduce the full actual-input measurement from the repository root:

```sh
sh benchmarks/run_actual.sh
```

`run_actual.sh` builds the binary, records environment/input identity, performs a separate streaming count, then invokes `measure.py`. The helper records the exact app command, wall time, child peak RSS, CPU time, and exit code. Logs are stored in:

- `environment.log`: Go/Python/platform and input size/hash.
- `count.stdout.log` and `count.stderr.log`: row/campaign count pass.
- `runtime.metrics.log`: timed application command and measurements.
- `actual.stdout.log` and `actual.stderr.log`: application streams.
- `output.sha256`: output hashes for regression comparison.
- `actual-results/`: output CSV files from the real dataset.

The first attempt using macOS `/usr/bin/time -l` ran the application but exited with a `sysctl kern.clockrate: Operation not permitted` error, so those timing-tool values are not used in the reported measurements. `getrusage` is used by the successful runs instead.

## CPU profile

The CLI accepts an optional `--cpuprofile FILE` flag. It starts `runtime/pprof` immediately before the application workflow and stops it after the workflow returns. Profiling is off by default. The profile file is synced and closed before the command returns; start, sync, and close errors are reported. The CLI rejects a profile path that aliases the input or either result file.

Capture a profile from the real dataset and compare generated CSVs with the unprofiled baseline:

```sh
GOCACHE=/private/tmp/go-build-cache sh benchmarks/run_profile.sh
```

The script builds the binary, runs one profiled pass, saves the profile at `benchmarks/ad_data.cpu.pprof`, captures app stdout/stderr and runtime/RSS metrics, hashes the outputs, and byte-compares both CSVs with `actual-results/`. Profiling adds overhead, so `profile.runtime.metrics.log` is not a baseline performance result.

Inspect the hottest functions and cumulative cost:

```sh
GOCACHE=/private/tmp/go-build-cache go tool pprof -top benchmarks/aggregator benchmarks/ad_data.cpu.pprof
GOCACHE=/private/tmp/go-build-cache go tool pprof -top -cum benchmarks/aggregator benchmarks/ad_data.cpu.pprof
GOCACHE=/private/tmp/go-build-cache go tool pprof -list='internal/csvdata.ReadEach' benchmarks/aggregator benchmarks/ad_data.cpu.pprof
```

The Go diagnostics documentation describes CPU profiles and the `go tool pprof` text and graphical reports ([Go diagnostics](https://go.dev/doc/diagnostics)). A call graph SVG can be generated when Graphviz `dot` is installed:

```sh
GOCACHE=/private/tmp/go-build-cache go tool pprof -svg benchmarks/aggregator benchmarks/ad_data.cpu.pprof > benchmarks/profile.svg
```

For an interactive pprof web UI, use `go tool pprof -http=:0 ...`; upstream pprof supports graph/flame graph views ([pprof README](https://github.com/google/pprof/blob/main/README.md)). This environment has no `dot`, so the SVG was not generated.

### Profile captured in this environment

The profile run completed with exit code 0 in 22.423 s and peak RSS of 12,206,080 bytes. Its output CSV hashes match the unprofiled baseline. The profile reported 22.72 s of samples over a 22.07 s profile duration:

| Sampled function/path | Flat samples | Cumulative samples |
|---|---:|---:|
| `syscall.rawsyscalln` | 82.92% | 82.92% |
| `encoding/csv.(*Reader).readRecord` | 0.13% | 83.41% |
| `bufio.(*Reader).ReadSlice` | 0% | 82.92% |
| `math/big.(*Rat).SetString` | 0% | 0.57% |

These are profile observations, not proof of a root cause. Samples are concentrated in the file-read/syscall path; the decimal parser and campaign aggregation have small sampled cost in this run. Further diagnosis should check the storage/cache environment and compare a repeated warm-cache profile before considering code changes. No optimization was made.

The exact profile, `top` report, stderr, runtime metrics, and output hashes are saved as `ad_data.cpu.pprof`, `profile.top.txt`, `profile.pprof.stderr.log`, `profile.runtime.metrics.log`, and `profile.output.sha256` in this directory.
