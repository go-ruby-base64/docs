# Usage & API

The public API lives at the module root (`github.com/go-ruby-base64/base64`). It
is **Ruby-shaped but Go-idiomatic**: `Encode64` / `Decode64` mirror MRI's
`Base64.encode64` / `Base64.decode64`, while the surface follows Go conventions —
operating on Go `string`s, an explicit `error` from the strict decoders, no global
state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-base64/base64`,
    bound into `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-base64/base64
```

## Worked example

```go
fmt.Printf("%q\n", base64.Encode64("hello world"))
// "aGVsbG8gd29ybGQ=\n"            — RFC-2045, trailing newline

fmt.Printf("%q\n", base64.StrictEncode64("hello"))
// "aGVsbG8="                      — no newline

fmt.Println(base64.Decode64("aGVs bG8=\n junk")) // hello — lenient

s, err := base64.StrictDecode64("not base64!!!")
fmt.Printf("%q %v\n", s, err)      // "" invalid base64

fmt.Println(base64.UrlsafeEncode64("\xfb\xff\xfe"))        // -__-
fmt.Println(base64.UrlsafeEncode64("ab", false))           // YWI  (unpadded)
```

## Shape

```go
// Encode64 — RFC 2045 (pack("m")): standard alphabet, \n every 60 chars + trailing \n.
func Encode64(s string) string

// StrictEncode64 — pack("m0"): standard alphabet, no newlines.
func StrictEncode64(s string) string

// Decode64 — lenient unpack1("m"): skips stray bytes, optional padding, padding stops.
func Decode64(s string) string

// StrictDecode64 — unpack1("m0"): returns ErrInvalid on any malformed input.
func StrictDecode64(s string) (string, error)

// UrlsafeEncode64 — -_ alphabet; padding defaults to true (pass false to strip it).
func UrlsafeEncode64(s string, padding ...bool) string

// UrlsafeDecode64 — -_ alphabet, RFC 4648; accepts padded or unpadded input.
func UrlsafeDecode64(s string) (string, error)

// ErrInvalid is returned by the strict decoders, mirroring MRI's ArgumentError.
var ErrInvalid = errors.New("invalid base64")
```

## Method semantics

- **`Encode64`** — RFC 2045 (`Array#pack("m")`): the standard `+/` alphabet, a
  newline every 60 output characters, and a trailing newline. The empty string
  encodes to the empty string (no newline), as in MRI.
- **`Decode64`** — lenient (`String#unpack1("m")`): non-alphabet bytes are
  skipped, padding is optional, a `=` on a 2- or 3-char partial quad finalises it
  and **stops** decoding (trailing-padding terminates the stream), a `=` on a quad
  boundary is ignored, and a lone orphaned final sextet is discarded.
- **`StrictEncode64` / `StrictDecode64`** — `pack("m0")` / `unpack1("m0")`: no
  newlines on encode; decode rejects any stray byte (including the embedded
  newlines `encoding/base64` would tolerate) with `ErrInvalid`.
- **`UrlsafeEncode64` / `UrlsafeDecode64`** — the `-_` alphabet (RFC 4648);
  padding is on by default and stripped when `padding=false`; decode accepts
  correctly-padded or unpadded input and rejects everything else.

## SIMD acceleration

`StrictEncode64`, `Encode64`, `StrictDecode64` and `UrlsafeDecode64` route the
standard-alphabet body through [`go-simd/base64`](https://github.com/go-simd/base64),
whose go-asmgen kernels cover the SIMD-capable 64-bit arches and fall back to
`encoding/base64` everywhere else, so the bytes are always identical to
`StdEncoding`. On amd64/arm64 the encode kernel measures markedly faster than the
scalar `encoding/base64` reference; the package ships `go test -bench` benchmarks
(`Benchmark*SIMD` vs `Benchmark*Scalar`) so the speedup can be measured on each
target — see [Performance](performance.md) for the methodology.

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** drives a wide
corpus — every length class, a line wrap, non-ASCII bytes, and malformed/strict
inputs — through both the system `ruby` and this library and compares the results
**byte-for-byte**. The oracle scripts `binmode` both stdin and stdout so Windows
text-mode never pollutes the bytes, gate on `RUBY_VERSION >= "4.0"`, and skip
themselves where `ruby` is absent (e.g. the qemu arch lanes), so the cross-arch
builds still validate the library.

## Relationship to Ruby

`go-ruby-base64/base64` is **standalone and reusable**, and is the Base64 backend
bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo`
as a native module — the same way
[go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-erb](https://github.com/go-ruby-erb) and
[go-ruby-yaml](https://github.com/go-ruby-yaml) are bound. Argument coercion and
raising `ArgumentError` are the host's job; the library returns a Go `error`
(`ErrInvalid`) the host maps to MRI's exception.
