package types

import "github.com/vegidio/go-sak/fetch"

type External interface {
	ExpandMedia(media Media, ignoreHost string, metadata *Metadata) Media
}

// Extractor defines the interface for extractors.
type Extractor interface {
	// Type returns the type of the extractor.
	Type() ExtractorType

	// SourceType determines the source type of the extractor.
	//
	// # Returns:
	//   - SourceType: the determined source type of the extractor.
	//   - error: non-nil if determination of the source type fails.
	SourceType() (SourceType, error)

	// QueryMedia queries media from the given URL with specified limit and extensions.
	//
	// # Parameters:
	//   - limit: maximum number of media items to return; extraction stops once this limit is reached.
	//   - extensions: list of file extensions (without a leading dot) to include in the results. If empty or nil, no
	//     extension-based filtering is applied.
	//   - deep: if true, performs a deep query on unknown URLs in an attempt to find extra media files.
	//
	// # Returns:
	//   - *Response: the Response object.
	//   - cancelFunc: a function to cancel the ongoing query.
	QueryMedia(limit int, extensions []string, deep bool) (*Response, func())

	// Fetch returns a fetch.Fetch instance configured to work with the extractor.
	//
	// # Parameters:
	//   - headers: the headers to use with the fetch.Fetch instance.
	Fetch(headers map[string]string) *fetch.Fetch
}
