# Roadmap

`go-ruby-base64/base64` is grown **test-first**, each capability
differential-tested against MRI rather than built in isolation. Ruby's `Base64`
module — the deterministic, interpreter-independent slice of `require "base64"` —
is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Encode64 | RFC 2045 (`pack("m")`): the standard `+/` alphabet, a newline every 60 output characters, and a trailing newline; the empty string encodes to empty. | **Done** |
| Decode64 (lenient) | `String#unpack1("m")`: non-alphabet bytes skipped, padding optional, a `=` on a partial quad finalises and **stops** decoding, a `=` on a quad boundary ignored, a lone orphaned final sextet discarded. | **Done** |
| Strict variants | `StrictEncode64` (`pack("m0")`) emits no newlines; `StrictDecode64` (`unpack1("m0")`) rejects any stray byte — including the newlines `encoding/base64` tolerates — with `ErrInvalid`. | **Done** |
| Url-safe variants | `UrlsafeEncode64` / `UrlsafeDecode64` on the `-_` alphabet (RFC 4648); padding on by default and stripped when `padding=false`; decode accepts padded or unpadded input and rejects everything else. | **Done** |
| SIMD acceleration | The hot standard-alphabet bodies route through [go-simd/base64](https://github.com/go-simd/base64)'s go-asmgen kernels (amd64 SSE/AVX2, arm64 NEON, ppc64le VSX, s390x vector; stdlib fallback elsewhere) — CGO=0 SIMD across all six 64-bit arches, bit-identical to `StdEncoding`. | **Done** |
| Differential oracle & coverage | A wide corpus — every length class, a line wrap, non-ASCII bytes, malformed/strict inputs — driven through the system `ruby` and asserted byte-equal; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic encode/decode
  algorithm over Go `string`s; it never runs arbitrary Ruby. Binding the methods
  to a live Ruby `Base64` module — argument coercion, raising `ArgumentError` — is
  the consumer's job, which is why `rbgo` binds this module rather than the
  reverse.
- **Go errors, not exceptions.** The strict decoders return a Go `error`
  (`ErrInvalid`); mapping it to MRI's `ArgumentError` is the host's job.
- **Reference is reference Ruby (MRI 4.0.5).** Byte-for-byte conformance targets
  MRI's behaviour, pinned by the differential oracle; the SIMD path stays
  bit-identical to `encoding/base64.StdEncoding`.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/SIMD split.
