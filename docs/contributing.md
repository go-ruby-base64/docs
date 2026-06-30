# Contributing

Contributions are welcome. `go-ruby-base64/base64` is built to a small set of
non-negotiable rules — they are what keep it pure-Go, correct, and
MRI-compatible. Please read these before opening a pull request.

## Hard rules

- **Build from source — no vendoring.** Everything compiles from source. Being
  able to compile from source is a guarantee of independence.
- **100% test coverage target, enforced in CI.** New code ships with tests, and
  coverage is a CI gate. Fill the error branches, not just the happy path.
- **All GitHub content in English.** Issues, pull requests, commits, comments,
  and discussions are English-only.
- **Differential testing against MRI.** Correctness is defined by reference Ruby.
  A corpus is run through both `ruby` and this library and the results are
  compared byte-for-byte — not approximated from memory.
- **Pure Go, cgo disabled.** The whole point is a single static binary with no C
  toolchain. Code must build with `CGO_ENABLED=0`. The SIMD path is go-asmgen
  Plan 9 assembly, not C — if a feature seems to need C, it needs a pure-Go path
  instead.
- **SIMD stays bit-identical.** The go-simd/base64 fast path must produce exactly
  the same bytes as `encoding/base64.StdEncoding` on every arch and fall back
  cleanly where no kernel exists.
- **A reusable library, not the interpreter.** This module implements the
  deterministic core. Argument coercion and raising `ArgumentError` belong in the
  consumer, not here.

## Workflow

1. Pick or open an issue describing the change.
2. Work test-first: add the differential / unit tests, then make them pass.
3. Run the full suite with coverage and confirm the gate is green:

    ```sh
    COVERPKG=$(go list ./... | paste -sd, -)
    go test -race -coverpkg="$COVERPKG" -coverprofile=cover.out ./...
    go tool cover -func=cover.out | tail -1   # 100.0%
    ```

4. Open a PR in English, referencing the issue.

## Where things live

The library is in
[`github.com/go-ruby-base64/base64`](https://github.com/go-ruby-base64/base64).
This documentation site is in
[`github.com/go-ruby-base64/docs`](https://github.com/go-ruby-base64/docs). Start
from the [Usage & API](api.md) page and the [Roadmap](roadmap.md) to find the
right place for your change.
