# Ad Performance Aggregator

A Go command-line program that streams an advertising CSV, aggregates metrics by campaign, and writes the highest-CTR and lowest-CPA rankings.

## Requirements

- Go 1.22 or later
- Input CSV with the required header and columns described below

The application uses only the Go standard library; there are no third-party Go dependencies.

## Build and run

Run directly from the repository root:

```sh
go run ./cmd/aggregator --input ad_data.csv --output results/
```

Build a standalone binary with:

```sh
go build -o aggregator ./cmd/aggregator
./aggregator --input ad_data.csv --output results/
```

To target a different OS or architecture, set `GOOS` and `GOARCH` to the desired Go target values before `go build`. For example, `GOOS=darwin GOARCH=amd64` builds for Intel-based Macs. A cross-compiled binary should be run and verified on its target environment.

The output directory is created when needed. Existing result files are replaced. The command writes `top10_ctr.csv` and `top10_cpa.csv` there. Both `--input` and `--output` are required; invalid arguments, unreadable input, malformed CSV, and processing or output errors are reported to stderr with a nonzero exit status.

To get the input data, place `ad_data.csv` in the repository root. If starting from the supplied archive, extract it before running:

```sh
unzip ad_data.csv.zip
```

## Input and output

The CSV must have this exact header and column order:

```csv
campaign_id,date,impressions,clicks,spend,conversions
```

Each row is validated. Dates must be valid `YYYY-MM-DD` dates; counters must fit signed 64-bit integers; spend must be a decimal number. A malformed header or row fails the run and includes row context in the error. Quoted CSV fields are parsed using Go's CSV reader.

The application sums impressions, clicks, spend, and conversions by `campaign_id`. It writes at most ten rows to each result file, with the columns:

```text
campaign_id,total_impressions,total_clicks,total_spend,total_conversions,CTR,CPA
```

Top CTR is sorted highest first; Top CPA is sorted lowest first and excludes campaigns with zero conversions. Campaigns with zero impressions are omitted from Top CTR because CTR is undefined. If CPA is undefined for a campaign included in Top CTR, its CPA field is empty. Ties are ordered by campaign ID ascending. Fewer than ten eligible campaigns produce fewer than ten rows.

Spend is parsed, summed, and compared as an exact decimal rational, avoiding binary floating-point accumulation error. Spend output preserves the exact terminating decimal value with at least two fractional digits. CPA ranking also uses exact rational comparison; CPA output has at least two and at most 18 fractional digits, with values requiring more precision rounded to 18 places. CTR is calculated and written to four decimal places.

## Design

The input is processed in a single streaming pass. Parsed records are passed directly to an in-memory map keyed by campaign ID; the program retains campaign totals rather than all input rows. Ranking happens after EOF by sorting the campaign aggregates, then the selected rows are written as CSV. Working memory therefore depends on the number of distinct campaigns and the largest input record, rather than the total number of rows.

Package responsibilities:

- `cmd/aggregator`: CLI flags, argument validation, error exit behavior, and optional CPU profiling.
- `internal/csvdata`: streaming CSV reading and field validation.
- `internal/aggregate`: campaign totals, derived metrics, and deterministic rankings.
- `internal/output`: CSV formatting and writing.
- `internal/app`: orchestration of reading, aggregation, ranking, and output.

The CSV, CLI, file, exact arithmetic, and profiling facilities come from the Go standard library, including `encoding/csv`, `flag`, `math/big`, and `runtime/pprof`. No external Go library is used.

## Tests and verification

Run the repository checks from its root:

```sh
go test ./...
go vet ./...
go build -o aggregator ./cmd/aggregator
```

A checked-in small fixture with expected outputs can be verified with:

```sh
sh benchmarks/run_fixture.sh
```

## Benchmark and memory

The real input benchmark was run on 2026-10-23 with the existing `ad_data.csv`; the measured application runs did not generate or decompress the input. Input identity was 1,043,304,870 bytes (about 995 MiB), SHA-256 `4e9c87cc0b408f043917f5f7ee0101ba09899ad79f154358e9a2da374cfd60eb`, 26,843,544 data rows, and 50 distinct campaigns. The environment was macOS Darwin 27.0.0, arm64, 8 logical CPUs, 16 GiB RAM, and Go `go1.27.3 darwin/arm64`.

Two unprofiled runs completed successfully in 21.302 s and 21.342 s wall time. Measured peak RSS was 12,959,744 bytes and 13,287,424 bytes, respectively (about 12.4–12.7 MiB). The second run measured 21.327 s user CPU and 0.674 s system CPU. Peak RSS and wall time were collected for the aggregator subprocess using Python `resource.getrusage` and `time.perf_counter`. These measurements include opening and reading the CSV, aggregation, ranking, and writing both result files; they exclude the separate counting and input-hash passes, build, and data extraction. They describe this machine and dataset, not a general performance guarantee.

Reproduce the actual-input measurement and inspect its command, logs, environment, and output hashes:

```sh
sh benchmarks/run_actual.sh
```

The script writes benchmark artifacts under `benchmarks/`. The recorded outputs from that run and the handoff outputs under `results/` are real results from the identified input; their hashes are recorded in `benchmarks/output.sha256` and can be compared with:

```sh
shasum -a 256 results/top10_ctr.csv results/top10_cpa.csv
```

Optional CPU profiling is available with `--cpuprofile FILE`; it is off by default. See [benchmarks/README.md](benchmarks/README.md) for profile collection, pprof commands, and the captured profile observations. No code optimization was made based on the profile.

## Handoff files

- `results/top10_ctr.csv` and `results/top10_cpa.csv`: outputs generated from the real `ad_data.csv` input.
- `benchmarks/`: fixture, measurement scripts, environment and runtime logs, output hashes, and CPU profile artifacts.
- `PROMPTS.md`: prompts used while developing the project.
