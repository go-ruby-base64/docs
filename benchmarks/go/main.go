// SPDX-License-Identifier: BSD-3-Clause
package main

import "github.com/go-ruby-base64/base64"

func main() {
	s := string(detBytes(3072)) // 3 KiB, identical to the Ruby workload
	enc := base64.StrictEncode64(s)

	bench("encode-3KiB", 2000, func() { sink = base64.StrictEncode64(s) })
	bench("decode-3KiB", 2000, func() { v, _ := base64.StrictDecode64(enc); sink = v })
}
