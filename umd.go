package umd

import (
	"fmt"

	"github.com/vegidio/umd/internal/extractors/coomer"
	"github.com/vegidio/umd/internal/extractors/fapello"
	"github.com/vegidio/umd/internal/extractors/imaglr"
	"github.com/vegidio/umd/internal/extractors/jpgfish"
	"github.com/vegidio/umd/internal/extractors/reddit"
	"github.com/vegidio/umd/internal/extractors/redgifs"
	"github.com/vegidio/umd/internal/extractors/saint"
	"github.com/vegidio/umd/internal/types"
)

// Umd represents a Universal Media Downloader instance.
type Umd struct {
	metadata types.Metadata
}

// New creates a new instance of Umd.
//
// # Parameters:
//   - metadata: A map containing metadata information.
//
// # Returns:
//   - Umd: A new instance of Umd.
func New(metadata types.Metadata) Umd {
	if metadata == nil {
		metadata = make(Metadata)
	}

	return Umd{metadata: metadata}
}

// FindExtractor attempts to find a suitable extractor for the given URL.
//
// # Parameters:
//   - url: The URL for which an extractor is to be found.
//
// # Returns:
//   - model.Extractor: The extractor instance if found.
//   - error: An error if no suitable extractor is found.
func (u Umd) FindExtractor(url string) (types.Extractor, error) {
	var extractor types.Extractor
	extractors := []func(string, types.Metadata, types.External) types.Extractor{
		coomer.New, fapello.New, imaglr.New, jpgfish.New, reddit.New, redgifs.New, saint.New,
	}

	for _, newExtractor := range extractors {
		if e := newExtractor(url, u.metadata, External{}); e != nil {
			extractor = e
			break
		}
	}

	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	return extractor, nil
}
