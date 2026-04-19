/*
Package sink provides ScopeSink, which injects a namespace or subsystem
identifier into every log entry before forwarding it downstream.

# Usage

	s := sink.NewScopeSink(inner, "payments", "scope")
	// every entry will carry Fields["scope"] = "payments"
	// unless the entry already sets that field.

The field name defaults to "scope" when an empty string is provided.
Existing field values on the entry take precedence over the sink's configured
scope, allowing per-entry overrides.
*/
package sink
