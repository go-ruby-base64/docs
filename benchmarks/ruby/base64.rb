# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "base64"
require_relative "_harness"

data = det_bytes(3072)         # 3 KiB, identical to the Go driver
enc  = Base64.strict_encode64(data)

bench("encode-3KiB", 2000) { Base64.strict_encode64(data) }
bench("decode-3KiB", 2000) { Base64.strict_decode64(enc) }
