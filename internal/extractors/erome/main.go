package erome

import (
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Erome struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "erome.com"):
		e := &Erome{}
		e.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Erome,
		}
		e.FetchMediaFn = e.fetchMedia
		e.SourceTypeFn = e.SourceType
		return e, nil
	}

	return nil, nil
}

var regexAlbum = regexp.MustCompile(`/a/([a-zA-Z0-9-_.]+)/?`)

func (e *Erome) SourceType() (types.SourceType, error) {

	var source types.SourceType

	switch {
	case regexAlbum.MatchString(e.Url):
		matches := regexAlbum.FindStringSubmatch(e.Url)
		id := matches[1]
		source = SourceAlbum{Id: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", e.Url)
	}

	e.Source = source
	return source, nil
}

func (e *Erome) DownloadHeaders() map[string]string {
	return map[string]string{
		"Referer": e.Url,
	}
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Erome)(nil)

// region - Private methods

func (e *Erome) fetchMedia(
	source types.SourceType,
	_ int,
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
			media, err := types.NewMedia(link, types.Erome, map[string]interface{}{
				"source":  strings.ToLower(sourceName),
				"name":    album.Id,
				"title":   album.Title,
				"created": album.Created,
			}, headers)
			if err != nil {
				continue
			}
			out <- media
		}
	}()

	return out
}

// endregion
