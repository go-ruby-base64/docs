# Performance

`go-ruby-base64/base64` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `Base64`. Its
hot standard-alphabet paths run on the SIMD kernels of
[go-simd/base64](https://github.com/go-simd/base64). This page records the
**methodology** for measuring it — both against the scalar `encoding/base64`
reference and against the reference Ruby runtimes.

!!! note "No numbers are published here yet"
    This page documents *how* the comparison is run, not a result table. Numbers
    are only added once they have been measured on the host described below and
    checked bit-identical to `StdEncoding` / byte-identical to MRI — never
    estimated or filled in from memory.

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
