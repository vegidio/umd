package umd

import (
	"fmt"

	"github.com/vegidio/umd/internal/extractors/coomer"
	"github.com/vegidio/umd/internal/extractors/cyberdrop"
	"github.com/vegidio/umd/internal/extractors/fapello"
	"github.com/vegidio/umd/internal/extractors/imaglr"
	"github.com/vegidio/umd/internal/extractors/jpgfish"
	"github.com/vegidio/umd/internal/extractors/reddit"
	"github.com/vegidio/umd/internal/extractors/redgifs"
	"github.com/vegidio/umd/internal/extractors/saint"
	"github.com/vegidio/umd/internal/extractors/simpcity"
	"github.com/vegidio/umd/internal/types"
)

// Umd represents a Universal Media Downloader instance.
type Umd struct {
	metadata types.Metadata
	headers  map[string]string
}

// New creates a new instance of Umd.
//
// # Parameters:
//   - metadata: A map containing metadata information.
//
// # Returns:
//   - Umd: A new instance of Umd.
func New() *Umd {
	return &Umd{metadata: make(Metadata), headers: nil}
}

// WithMetadata sets the metadata for the Umd instance.
//
// # Parameters:
//   - metadata: A map containing metadata information.
//
// # Returns:
//   - Umd: An updated instance of Umd.
func (u *Umd) WithMetadata(metadata types.Metadata) *Umd {
	u.metadata = metadata
	return u
}

// FindExtractor attempts to find a suitable extractor for the given URL.
//
// # Parameters:
//   - url: The URL for which an extractor is to be found.
//
// # Returns:
//   - model.Extractor: The extractor instance if found.
//   - error: An error if no suitable extractor is found.
func (u *Umd) FindExtractor(url string) (types.Extractor, error) {
	var extractor types.Extractor
	extractors := []func(string, types.Metadata, types.External) (types.Extractor, error){
		coomer.New,
		cyberdrop.New,
		fapello.New,
		imaglr.New,
		jpgfish.New,
		reddit.New,
		redgifs.New,
		saint.New,
		simpcity.New,
	}

	for _, newExtractor := range extractors {
		if ext, err := newExtractor(url, u.metadata, External{}); ext != nil {
			if err != nil {
				return nil, err
			}

			extractor = ext
			break
		}
	}

	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	return extractor, nil
}
