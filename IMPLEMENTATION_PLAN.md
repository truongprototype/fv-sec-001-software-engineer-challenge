# Implementation plan — Ad Performance Aggregator

This plan is based on the complete project `README.md`. No implementation code or result data is included here. Statements are marked **Required** when the README makes them part of the challenge, and **Proposal** when they are implementation choices suggested by this plan.

## 1. Requirements extracted from README

### Functional requirements

- **Required:** Build a console/CLI application that reads the supplied CSV dataset (about 1 GB) and produces aggregated advertising analytics.
- **Required:** Input columns are `campaign_id` (string), `date` (string, `YYYY-MM-DD`), `impressions` (integer), `clicks` (integer), `spend` (float, USD), and `conversions` (integer).
- **Required:** Aggregate records by `campaign_id`, summing impressions, clicks, spend, and conversions.
- **Required:** Calculate CTR as `total_clicks / total_impressions` and CPA as `total_spend / total_conversions`.
- **Required:** When conversions are zero, CPA may be ignored or represented as null. Campaigns with zero conversions must be excluded from the lowest-CPA list.
- **Required:** Produce `top10_ctr.csv` containing the ten campaigns with highest CTR, and `top10_cpa.csv` containing the ten campaigns with lowest CPA. Rows include campaign ID, four aggregate totals, CTR, and CPA as shown in the README.
- **Required:** Provide a CLI accepting an input path and output directory (the README gives `--input` and `--output` as an example).

### Technical and quality requirements

- **Required:** Process a roughly 1 GB CSV efficiently, with good performance and memory optimization.
- **Required:** Handle missing files, malformed rows, and edge cases gracefully.
- **Required:** Include tests that verify correctness.
- **Required:** Keep code clean and readable, use meaningful names, consistent style, and avoid dead/commented-out code.
- **Required:** Document non-obvious decisions in the README.

### Handoff / submission requirements

- **Required:** Submit a GitHub repository link by email to `backoffice@flinters.vn`.
- **Required:** Repository includes source code and the two output files `top10_ctr.csv` and `top10_cpa.csv`.
- **Required:** Repository README explains setup, invocation, libraries used, processing time for the 1 GB file, and peak memory if measured.
- **Optional (README):** Dockerfile and benchmark logs.
- **Conditional / optional (README):** If AI coding assistants are used, include root `PROMPTS.md`, with prompts pasted as-is. The section calls this not mandatory but highly valued, despite saying “if used.”

## 2. Edge cases and unspecified decisions

These points are not fully specified by the README and should be settled/documented before implementation. Items below are questions, not added requirements.

1. **Zero impressions:** CTR division is undefined. Should such campaigns be excluded from Top CTR, assigned null/zero CTR, or treated another way?
2. **Zero conversions / CPA field:** README permits ignoring or null CPA for zero conversions, and excludes those campaigns from Top CPA. It does not specify how CPA should appear in the Top CTR CSV for such campaigns.
3. **Malformed records:** The README asks for graceful handling but does not say whether invalid rows are skipped, fail the run, or are collected and reported; nor what constitutes malformed data.
4. **Header and column order:** Header presence, exact header matching, column order, extra/missing columns, BOM, quoting, and blank lines are not specified.
5. **Numeric validity:** It does not define whether negative counters/spend, fractional counters, NaN/infinity, empty numeric fields, or values beyond Go integer ranges are allowed.
6. **Date validity and use:** Date format is stated, but calendar-date validation and date-based aggregation/filtering are not requested. The date appears to be input metadata only; whether invalid dates reject a row is unspecified.
7. **Money precision:** `spend` is described as float USD, but precision/rounding policy and output decimal scale are unspecified. Examples show two decimal places.
8. **Tie-breaking:** No ordering is prescribed for equal CTR or CPA, including ties at position 10. Output should be deterministic, but the chosen secondary key is unspecified.
9. **Fewer than ten qualifying campaigns:** Whether to output all available campaigns or fail/pad is unspecified.
10. **Duplicate campaign IDs / empty IDs:** Grouping implies repeated campaign IDs combine, but whitespace normalization and empty IDs are unspecified.
11. **Overflow:** Maximum rows and per-field values are unspecified; accumulated totals may overflow fixed-width integers.
12. **CSV output details:** Header is implied by the sample table; exact float formatting, null representation, line endings, and ordering are unspecified.
13. **Output lifecycle:** Behavior if output directory does not exist, files already exist, or writing fails partway is unspecified.
14. **CLI contract:** Defaults, flags beyond input/output, exit codes, diagnostics, and whether output names are fixed or configurable are unspecified.
15. **Input delivery:** README describes downloading/unzipping the dataset, but does not require bundling the large input file in the repository.
16. **Benchmark environment:** Required README metric is processing time for the 1 GB file; peak memory is requested only if measured. Hardware, OS, Go version, cache state, and timing method are not defined.

## 3. Minimal Go architecture (proposal)

**Proposal:** Keep a small standard-library-only module with narrow components and no long-lived per-row storage:

- `cmd/aggregator/main.go`: parse CLI flags, call application workflow, print concise errors, set nonzero exit status on failure.
- `internal/aggregate`: `Campaign` aggregate type and pure aggregation/ranking functions. Use a map keyed by `campaign_id`; retain only one aggregate per distinct campaign.
- `internal/csvinput`: read records from `encoding/csv.Reader` sequentially, validate field count and parse fields, and return one row at a time (or invoke a callback). Avoid `ReadAll`.
- `internal/output`: serialize the selected rows via `encoding/csv.Writer` to the two fixed output files.
- `internal/app` (optional if it keeps `main` simple): orchestrate streaming input, aggregation, selection, and output; inject reader/writer dependencies in tests.

The streaming pass uses memory proportional to the number of distinct campaign IDs, not the 1 GB file size. If the distinct campaign count itself is too large for memory, the README does not require an external-sort/spill design; treat that as a scale assumption to verify with the actual dataset before expanding scope. Use buffered standard-library CSV I/O; choose an explicit `FieldsPerRecord` policy after resolving schema/header handling. Do not add concurrency unless measurement shows a need, because aggregation is shared-state work and one sequential pass is the simplest baseline.

Suggested errors should carry row number and field context. Keep parsing, aggregation, and ranking independently testable without filesystem or CLI setup.

## 4. Aggregation and Top 10 algorithms (proposal)

1. Read and parse each CSV record once.
2. For each accepted record, find/create `campaign_id` aggregate and add impressions, clicks, conversions, and spend.
3. After EOF, compute CTR and CPA from aggregate totals; explicitly represent undefined ratios (zero denominator) rather than allowing NaN/Inf into sorting or CSV.
4. Rank Top CTR descending and Top CPA ascending. With `C` distinct campaigns, full sort costs `O(C log C)` time and `O(C)` memory. It is simple, deterministic, and likely adequate if campaign cardinality is moderate.
5. If `C` is large enough to make full sorting material, maintain a size-10 heap per ranking: `O(C log 10)` selection time and `O(10)` extra memory, then sort each selected set by the final comparator. The campaign aggregate map still costs `O(C)`.
6. Select Top CPA only from campaigns with conversions greater than zero. Top CTR handling of zero impressions needs a policy decision.

**Deterministic tie proposal:** Compare primary metric first; on equality, compare `campaign_id` lexicographically ascending. Apply the same secondary order at the 10th-place cutoff and when ordering final output. This is a proposal, not a README requirement.

## 5. Data and arithmetic policies (proposal)

- **Money:** Prefer integer cents if the input contract guarantees at most two fractional decimal digits; otherwise parse decimal text exactly into a scaled integer with a documented scale, or use a decimal type. The README says `float` and does not guarantee scale, so do not silently round to cents without confirming the data contract. Avoid binary floating-point accumulation when exact cent totals are expected. Format spend and CPA consistently (examples use two decimals); decide and document rounding.
- **CTR / CPA:** Use floating-point ratios for ranking and output if the data contract accepts approximate ratios; guard denominator zero before division. Do not serialize `NaN` or infinity. For an undefined CPA in Top CTR rows, a practical proposal is an empty CSV field (null-like); confirm desired representation. Top CPA excludes zero-conversion campaigns as explicitly required.
- **Zero impressions:** Proposed policy is to omit such campaigns from Top CTR because CTR is undefined; document and test once accepted. Do not silently treat as 0 without a decision.
- **Malformed data:** Proposed default is fail-fast on malformed header/schema or a structurally invalid record, with row number and reason; optionally report the first error only. Skipping rows risks silently incorrect totals. If product expectations prefer best-effort processing, define skipped-row counts and diagnostics before implementation.
- **Integer parsing and overflow:** Parse counters as base-10 signed 64-bit integers, reject parse failures, and use checked addition or detect overflow. Whether negative values are valid remains an open data rule; do not reject them merely by assumption.
- **Date / identifiers:** Preserve campaign ID text as read unless normalization is explicitly specified. Validate date shape (and optionally calendar validity) only according to a decided malformed-row policy; date does not affect grouping/ranking.
- **CSV correctness:** Use `encoding/csv` for quoted fields and `encoding/csv.Writer` for output, rather than manual comma splitting/concatenation.

## 6. Reproducible test and benchmark plan

The README requires tests; it does not prescribe particular cases or a benchmark procedure. The following are proposals. These are plans only; no test data or result files are created by this document.

### Unit / integration tests (proposal)

- Aggregation combines multiple dates/rows for the same campaign and keeps different campaigns separate.
- CTR / CPA calculations use aggregate totals, not averages of per-row ratios.
- Ranking returns correct highest CTR and lowest eligible CPA, includes no more than 10 rows, and handles fewer than 10 candidates.
- Zero conversions are absent from Top CPA; their Top CTR CPA cell follows the chosen null policy.
- Zero impressions follow the chosen CTR policy.
- Tied metrics, including a tie at the cutoff, follow the deterministic comparator.
- CSV parsing covers header, quoted fields, CRLF, blank lines, wrong field counts, invalid integer/decimal, invalid date under the chosen rule, and empty campaign ID under the chosen rule.
- Numeric boundary tests cover accumulated totals and overflow handling.
- CLI integration verifies missing input, output directory behavior, output headers/filenames, valid small input, and write errors.
- Output round-trip/expected-string assertions verify CSV escaping and stable numeric formatting.

Use small inline fixture strings in tests rather than checked-in/generated result datasets unless a fixture is specifically needed.

### Benchmark (proposal)

- Add Go benchmarks around streaming parse + aggregation and around ranking separately, using a deterministic fixture generated in-memory by the benchmark itself; do not commit a synthetic result dataset.
- For the end-to-end 1 GB run, record exact input byte size, SHA-256, Go version, OS/architecture, CPU/RAM, command, wall-clock elapsed time, and peak RSS if available. Use the same input and environment for comparisons; disclose filesystem-cache conditions.
- Run with a clean output directory and capture benchmark logs as optional submission material. Keep generated output results separate from benchmark fixtures.
- For memory behavior, use `/usr/bin/time -v` where available (or a platform equivalent) and report peak resident set size. State that aggregation memory scales with distinct campaign count.
- Run repeat measurements (e.g. three runs) and report each plus median; avoid claiming a general performance number from one run.

## 7. Decisions to resolve before implementation

Recommended minimal decisions: header/schema policy; fail-fast versus skip malformed rows; counter/spend validity and precision; zero-impression CTR handling; null representation for CPA in Top CTR; deterministic tie key; fewer-than-ten behavior; output overwrite behavior. Record the chosen rules in the project README as required documentation of non-obvious decisions.
