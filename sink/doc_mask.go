/*
Package sink provides MaskSink, a sink decorator that partially obscures
string field values before forwarding log entries downstream.

# Overview

MaskSink is useful when logs may contain sensitive identifiers such as payment
card numbers, API tokens, or national IDs where you want to retain a recognisable
prefix and/or suffix for debugging while hiding the sensitive middle portion.

# Usage

	s := sink.NewMaskSink(inner,
		sink.MaskOption{
			Field:      "card_number",
			KeepPrefix: 4,
			KeepSuffix: 4,
			Mask:       "****",
		},
	)

If the value is shorter than KeepPrefix+KeepSuffix it is left unchanged.
The Mask string defaults to "****" when empty.
*/
package sink
