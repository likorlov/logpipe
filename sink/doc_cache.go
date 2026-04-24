/*
Package sink provides CacheSink, a sink wrapper that short-circuits repeated
writes of identical log entries within a configurable TTL.

# Overview

CacheSink is useful when a hot code path emits the same structured log message
at very high frequency and forwarding every occurrence to an expensive inner
sink (e.g. a webhook or file) would be wasteful. The cache key is the value of
the "message" field; entries without that field are always forwarded.

# Usage

	inner := sink.NewConsoleSink(os.Stdout, false)
	cached := sink.NewCacheSink(inner, 5*time.Second)
	defer cached.Close()

	logger := logpipe.New(logpipe.InfoLevel, cached)
	logger.Info("disk usage high", logpipe.Fields{"pct": 92})
	// Subsequent writes with message "disk usage high" within 5 s are dropped.

# Behaviour

  - The first occurrence of a message within the TTL window is forwarded and
    its result (nil or error) is stored in the cache.
  - Subsequent occurrences within the TTL return the cached result immediately
    without calling the inner sink.
  - After the TTL elapses the next write is forwarded and refreshes the cache.
  - Entries with no "message" field bypass the cache entirely.
*/
package sink
