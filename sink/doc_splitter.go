/*
Package sink — SplitterSink

SplitterSink routes each log entry to one of several named child sinks based on
a user-supplied SplitterFunc. This is useful when you want to fan entries out to
different destinations depending on a field value such as level, service, or
region.

# Routing

The SplitterFunc receives the entry and returns a string key. If the key matches
a registered route the entry is forwarded there. If no match is found the entry
is sent to the sink registered under defaultKey. If defaultKey is empty or
unregistered the entry is silently dropped.

# Example

	splitter := sink.NewSplitterSink(
	    func(e logpipe.Entry) string {
	        if lvl, ok := e["level"].(string); ok {
	            return lvl
	        }
	        return "default"
	    },
	    "default",
	    map[string]logpipe.Sink{
	        "error":   sink.NewFileSink("/var/log/errors.log"),
	        "default": sink.NewConsoleSink(false),
	    },
	)
*/
package sink
