/*
Package sink provides BatchSink, which collects log entries into an
in-memory batch and forwards them together to a caller-supplied flush
function.

# Flush triggers

A flush is triggered by whichever condition occurs first:
  - The batch reaches the configured maximum size (maxSize).
  - The flush interval timer fires.
  - Close is called (remaining entries are flushed synchronously).

# Example

	flushFn := func(batch []logpipe.Entry) error {
		for _, e := range batch {
			fmt.Println(e.Message)
		}
		return nil
	}

	s := sink.NewBatchSink(50, 5*time.Second, flushFn)
	defer s.Close()

	logger := logpipe.New(s)
	logger.Info("hello batch")
*/
package sink
