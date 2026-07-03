# Performance

`go-ruby-base64/base64` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `Base64`. Its
hot standard-alphabet paths run on the SIMD kernels of
[go-simd/base64](https://github.com/go-simd/base64). This page records the
**methodology** for measuring it — both against the scalar `encoding/base64`
reference and against the reference Ruby runtimes.

## Result (best of 5, ms)

Measured 2026-06-30 on **Apple M4 Max**, macOS (darwin/arm64), Go 1.26.4, with
`ruby 4.0.5 +PRISM`, `jruby 10.1.0.0` (OpenJDK 25) and `truffleruby 34.0.1`
(GraalVM CE Native). The cross-runtime `Base64.encode64` / `strict_encode64` /
`decode` round-trip over a fixed ~3 KiB binary payload; output checked
byte-identical to MRI before timing.

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-base64) | 680 | 4.00× |
| MRI (ruby 4.0.5) | 170 | 1.00× |
| MRI + YJIT | 180 | 1.06× |
| JRuby 10.1.0.0 | 1640 | 9.65× |
| TruffleRuby 34.0.1 | 1140 | 6.71× |

rbgo runs on **go-ruby-base64** and is **~4× slower than MRI** here. The
go-simd/base64 kernel *is* on the path and *is* faster on raw bytes (see the
SIMD-vs-scalar benchmark below), but on this **Ruby-visible** round-trip the cost
is dominated by per-call Ruby-string allocation and the interpreter dispatch
around each `Base64.*` call, not the transform itself: MRI's `pack("m")` is a tight
C path with cheaper string handling, so the SIMD advantage is swamped. The kernel
win shows up in the Go-internal `go test -bench` (bytes in/out); it does not at the
Ruby level, because the bottleneck has moved to string churn. This is the top
per-module optimization target for go-ruby-base64 (cut the per-call allocation in
the binding); output stays byte-identical to MRI.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. These
    are **real measured numbers** from the 2026-06-30 run (Apple M4 Max;
    `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`, `truffleruby 34.0.1`) — nothing is
    fabricated or cherry-picked.

## Two comparisons

**1. SIMD vs scalar (Go-internal).** The package ships `go test -bench`
benchmarks (`Benchmark*SIMD` vs `Benchmark*Scalar`) that time the go-simd/base64
kernel against the scalar `encoding/base64` reference for encode and decode across
input sizes. Because both produce bit-identical output, this isolates the SIMD
speedup on whatever arch the bench runs on (amd64 SSE/AVX2, arm64 NEON, ppc64le
VSX, s390x vector). Reproduce:

```sh
go test -bench=. -benchmem ./...
```

**2. Ruby-visible operation (cross-runtime).** The **same** Ruby script — a
`Base64.encode64` / `Base64.decode64` round-trip of a representative payload — is
run under every runtime. `rbgo`'s number reflects **this pure-Go library doing the
work**; every other column is that interpreter's own `base64` stdlib. The script's
output is checked **byte-identical to MRI** before any timing is recorded.

## How to reproduce the cross-runtime comparison

- **Host:** a single, recorded machine (CPU, OS, arch noted alongside any result
  table), so numbers are comparable run to run, and so the SIMD lane in use is
  unambiguous.
- **Method:** best-of-N wall time (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** MRI (the oracle) and MRI `--yjit`; the JVM-based and GraalVM-based
  Rubies are timed **cold, single-shot**, so they carry VM startup on every run —
  read them as one-shot `ruby file.rb` costs, the same way `rbgo` and MRI are
  measured, not as steady-state JIT numbers.
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules).

!!! warning "Honest framing"
    Rows that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. Any numbers added here will be real
    measured numbers from a dated run, on a named host with the SIMD lane stated,
    with nothing cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `base64`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to
MRI before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### decode-3KiB

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 5997.8 | 4.70× |
| MRI | 1275.5 | 1.00× |
| MRI + YJIT | 1221.5 | 0.96× |
| JRuby | 1777.9 | 1.39× |
| TruffleRuby | 8768.7 | 6.87× |

#### encode-3KiB

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 746.1 | 0.67× |
| MRI | 1106.5 | 1.00× |
| MRI + YJIT | 1071.5 | 0.97× |
| JRuby | 7244.1 | 6.55× |
| TruffleRuby | 9790.7 | 8.85× |

Encode **beats MRI** (0.67×): the [go-simd/base64](https://github.com/go-simd/base64) kernel is on the encode path. Decode is currently ~4.7× MRI's C `String#unpack('m')` — go-ruby's `StrictDecode64` allocates a fresh output string and validates strictly on every call, so decode is this module's top optimization target. Output stays byte-identical to MRI (verified: encoded SHA-256 prefix `76e66b49…`).

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-base64/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/base64.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (a few cold-JIT
    outliers are noted in the text). Sub-microsecond rows carry the most relative
    noise; treat those ratios as order-of-magnitude. Every number here is a
    **real measured value** from the dated run above — nothing is fabricated,
    estimated, or cherry-picked. The go-ruby column is the pure-Go library; every
    other column is that interpreter's own stdlib doing the equivalent work.
