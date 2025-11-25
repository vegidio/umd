package erome

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Erome struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "erome.com"):
		return &Erome{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (e *Erome) Type() types.ExtractorType {
	return types.Erome
}

func (e *Erome) SourceType() (types.SourceType, error) {
	regexAlbum := regexp.MustCompile(`/a/([a-zA-Z0-9-_.]+)/?`)

	var source types.SourceType

	switch {
	case regexAlbum.MatchString(e.url):
		matches := regexAlbum.FindStringSubmatch(e.url)
		id := matches[1]
		source = SourceAlbum{Id: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", e.url)
	}

	e.source = source
	return source, nil
}

func (e *Erome) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if e.responseMetadata == nil {
		e.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       e.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Erome,
		Metadata:  e.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if e.source == nil {
			e.source, err = e.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := e.fetchMedia(e.source, limit, extensions, deep)

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

				// Limiting the number of results
				if utils.MergeMedia(&response.Media, result.Data) >= limit {
					response.Media = response.Media[:limit]
					return
				}
			}
		}
	}()

	return response, stop
}

func (e *Erome) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Erome)(nil)

// region - Private methods

func (e *Erome) fetchMedia(
	source types.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var posts <-chan saktypes.Result[Album]

		switch s := source.(type) {
		case SourceAlbum:
			posts = e.fetchAlbum(s)
		}

		for post := range posts {
			if post.Err != nil {
				out <- saktypes.Result[types.Media]{Err: post.Err}
				return
			}

			media := e.dataToMedia(post.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (e *Erome) fetchAlbum(source SourceAlbum) <-chan saktypes.Result[Album] {
	result := make(chan saktypes.Result[Album])

	go func() {
		defer close(result)
		album, err := getAlbum(source.Id)

		if err != nil {
			result <- saktypes.Result[Album]{Err: err}
		} else {
			result <- saktypes.Result[Album]{Data: *album}
		}
	}()

	return result
}

func (e *Erome) dataToMedia(album Album, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := e.DownloadHeaders()

	go func() {
		defer close(out)

		for _, link := range album.Links {
			out <- types.NewMedia(link, types.Erome, map[string]interface{}{
				"id":      album.Id,
				"name":    album.Name,
				"source":  strings.ToLower(sourceName),
				"created": album.Created,
			}, headers)
		}
	}()

	return out
}

// endregion
