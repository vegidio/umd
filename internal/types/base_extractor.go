package types

import (
	"context"

	saktypes "github.com/vegidio/go-sak/types"
)

// FetchMediaFunc is the signature for an extractor's media fetching function.
type FetchMediaFunc func(source SourceType, limit int, extensions []string, deep bool) <-chan saktypes.Result[Media]

// BaseExtractor provides the common fields and QueryMedia/Type implementations shared by all extractors.
type BaseExtractor struct {
	Metadata         Metadata
	Url              string
	Source           SourceType
	ResponseMetadata Metadata
	External         External
	ExtType          ExtractorType
	FetchMediaFn     FetchMediaFunc
	SourceTypeFn     func() (SourceType, error)
}

func (b *BaseExtractor) Type() ExtractorType {
	return b.ExtType
}

func (b *BaseExtractor) QueryMedia(limit int, extensions []string, deep bool) (*Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if b.ResponseMetadata == nil {
		b.ResponseMetadata = make(Metadata)
	}

	response := &Response{
		Url:       b.Url,
		Media:     make([]Media, 0),
		Extractor: b.ExtType,
		Metadata:  b.ResponseMetadata,
		Done:      make(chan error, 1),
	}

	seen := make(map[string]struct{})

	go func() {
		defer close(response.Done)

		if b.Source == nil {
			b.Source, err = b.SourceTypeFn()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := b.FetchMediaFn(b.Source, limit, extensions, deep)

		for {
			select {
			case <-ctx.Done():
				return

			case result, ok := <-mediaCh:
				if !ok {
					return
				}

				if result.Err != nil {
					response.Done <- result.Err
					return
				}

				// Deduplicate and limit results
				response.Mu.Lock()
				if _, exists := seen[result.Data.Url]; !exists {
					seen[result.Data.Url] = struct{}{}
					response.Media = append(response.Media, result.Data)
				}
				count := len(response.Media)
				if count >= limit {
					response.Media = response.Media[:limit]
					response.Mu.Unlock()
					return
				}
				response.Mu.Unlock()
			}
		}
	}()

	return response, stop
}
