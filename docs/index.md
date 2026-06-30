# go-ruby-base64 documentation

**Ruby's `Base64` module in pure Go — MRI-compatible, SIMD-accelerated, no cgo.**

`go-ruby-base64/base64` is a faithful, pure-Go (zero cgo) reimplementation of
Ruby's standard-library
[`Base64`](https://docs.ruby-lang.org/en/master/Base64.html) module — the
deterministic, interpreter-independent core of MRI 4.0.5's `require "base64"`,
matching reference Ruby byte-for-byte. The module path is
`github.com/go-ruby-base64/base64`.

The hot standard-alphabet encode/decode paths run on the SIMD kernels of
[**go-simd/base64**](https://github.com/go-simd/base64) (go-asmgen: amd64
SSE/AVX2, arm64 NEON, ppc64le VSX, s390x vector; stdlib fallback elsewhere), so
the output stays bit-identical to `encoding/base64.StdEncoding` while going faster
on the supported arches — a **CGO=0 SIMD path on all six of Go's 64-bit
targets**. Only the MRI-specific framing — the 60-column wrap, the lenient
`unpack("m")` state machine, and the url-safe padding rules — is hand-written.

It was **extracted into a reusable standalone library**: importable by any Go
program with no Ruby runtime, and bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a native
module — just like [go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-erb](https://github.com/go-ruby-erb) and
[go-ruby-yaml](https://github.com/go-ruby-yaml).

!!! success "Status: complete — Base64 byte-exact"
    **`Encode64`** (RFC 2045, newline every 60 chars + trailing newline), lenient **`Decode64`** (skips stray bytes, optional padding, padding stops the stream), strict **`StrictEncode64`** / **`StrictDecode64`** (no newlines; `ErrInvalid` on any malformed byte), and **`UrlsafeEncode64`** / **`UrlsafeDecode64`** (the `-_` alphabet, optional padding). The hot paths route through go-simd/base64 and stay bit-identical to `encoding/base64.StdEncoding`. Validated by a **differential oracle** against the system `ruby` byte-for-byte, at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
base64.Encode64("hello world")          // "aGVsbG8gd29ybGQ=\n"  — RFC-2045, trailing \n
base64.StrictEncode64("hello")          // "aGVsbG8="            — no newline
base64.Decode64("aGVs bG8=\n junk")     // "hello"              — lenient

s, err := base64.StrictDecode64("not base64!!!") // "", ErrInvalid
base64.UrlsafeEncode64("\xfb\xff\xfe")  // "-__-"
base64.UrlsafeEncode64("ab", false)     // "YWI"                — unpadded
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`base64`](https://github.com/go-ruby-base64/base64) | the library — Ruby's `Base64` in pure Go, SIMD-accelerated |
| [`docs`](https://github.com/go-ruby-base64/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-base64.github.io`](https://github.com/go-ruby-base64/go-ruby-base64.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-base64/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **SIMD-accelerated, bit-identical.** The standard-alphabet body runs on
  go-simd/base64's go-asmgen kernels across all six 64-bit Go arches, falling back
  to `encoding/base64` everywhere else — always identical to `StdEncoding`.
- **MRI byte-exact.** Output matches reference Ruby exactly, validated by a
  differential oracle against the `ruby` binary.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why Base64 is deterministic enough to live as a
  standalone, interpreter-independent SIMD-accelerated Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-base64/base64](https://github.com/go-ruby-base64/base64).
