package utils

import "context"

// Send delivers v on out, unless ctx is cancelled first. It reports whether the value was delivered;
// a false result means the consumer is gone and the caller should stop producing.
//
// Every extractor produces media through a chain of goroutines feeding unbuffered channels. A plain
// `out <- v` there blocks forever as soon as the consumer walks away - which it does on every query
// that hits its limit or is cancelled - leaking the whole chain.
func Send[T any](ctx context.Context, out chan<- T, v T) bool {
	select {
	case out <- v:
		return true
	case <-ctx.Done():
		return false
	}
}
