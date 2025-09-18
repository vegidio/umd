package types

import (
	"fmt"
	"time"
)

// Response represents a response from a service.
type Response struct {
	// Url is the URL from which the response was obtained.
	Url string

	// Media is a list of Media objects associated with the response.
	Media []Media

	// Extractor is the type of extractor used to obtain the response.
	Extractor ExtractorType

	// Metadata contains additional metadata about the response.
	Metadata Metadata

	// Done is a channel used to signal when the media query is complete.
	Done chan error
}

// Error waits for the query to finish and returns any error that occurred during the process.
func (r *Response) Error() error {
	err := <-r.Done
	return err
}

// Track monitors changes in the number of Media items and invokes the callback with queried and total counts. The
// callback receives the number of Media items queried and the total number of Media items.
//
// # Parameters:
//   - callback: A function that takes two arguments: current queried media (int), total number of queried media (int).
//
// # Returns:
//   - An error if one occurred during the query process.
func (r *Response) Track(callback func(queried, total int)) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	oldValue := 0

	for {
		select {
		case <-ticker.C:
			total := len(r.Media)
			if total != oldValue {
				callback(total-oldValue, total)
				oldValue = total
			}

		case <-r.Done:
			total := len(r.Media)
			if total != oldValue {
				callback(total-oldValue, total)
			}
			return r.Error()
		}
	}
}

func (r *Response) String() string {
	return fmt.Sprintf("{Url: %s, Media: %v, Extractor: %s, Metadata: %v}",
		r.Url, r.Media, r.Extractor, r.Metadata)
}
