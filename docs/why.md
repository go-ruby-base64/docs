# Why pure Go

`go-ruby-base64/base64` reimplements Ruby's `Base64` module in **pure Go, with cgo
disabled**. Base64 encoding and decoding are **fully deterministic and
interpreter-independent**: given the input bytes, the result is a pure function of
those bytes — no live binding, no evaluation of arbitrary Ruby. That is exactly
the part that can — and should — live as a standalone Go library, separate from
the interpreter.

## SIMD without cgo

The hot standard-alphabet encode/decode bodies run on
[**go-simd/base64**](https://github.com/go-simd/base64), whose
[go-asmgen](https://github.com/go-asmgen) kernels cover the SIMD-capable 64-bit
arches — amd64 SSE/AVX2, arm64 NEON, ppc64le VSX, s390x vector — and fall back to
`encoding/base64` everywhere else. Because go-asmgen emits **Plan 9 assembly**,
the SIMD path needs **no C and no cgo**: it cross-compiles and links into a single
static binary like the rest of the package, and the output is always
**bit-identical to `encoding/base64.StdEncoding`**. This is the same CGO=0 SIMD
approach the go-simd ecosystem applies across all six of Go's 64-bit targets.

Only the MRI-specific framing is hand-written here: the RFC-2045 60-column wrap,
the lenient `unpack("m")` quad/padding state machine, and the url-safe padding
rules. Everything else delegates to the SIMD kernel.

## Extracted from rbgo, reusable by anyone

This library is the Base64 backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)'s `rbgo`, but it has
been **extracted into a reusable standalone library** so that:

- any Go program can import `github.com/go-ruby-base64/base64` directly, with no
  Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native
  module (the same pattern as [go-ruby-regexp](https://github.com/go-ruby-regexp),
  [go-ruby-erb](https://github.com/go-ruby-erb) and
  [go-ruby-yaml](https://github.com/go-ruby-yaml)), rather than this module
  depending on the interpreter;
- the behaviour is pinned by a **differential oracle** against the system `ruby`,
  independent of any one consumer.

## Why pure Go matters here

Because the library is CGO-free and dependency-free, it:

- cross-compiles to every Go target with no C toolchain, and links into a single
  static binary — SIMD path included;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on
  `PATH`, while the cross-arch lanes (where `ruby` is absent) still validate the
  library.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is
in scope.
