package fapello

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Fapello struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "fapello.com"):
		f := &Fapello{}
		f.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Fapello,
		}
		f.FetchMediaFn = f.fetchMedia
		f.SourceTypeFn = f.SourceType
		return f, nil
	}

	return nil, nil
}

var (
	regexPost   = regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/(\d+)`)
	regexModel  = regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/?`)
	regexPostId = regexp.MustCompile(`/(\d+)/?$`)
)

func (f *Fapello) SourceType() (types.SourceType, error) {

	var source types.SourceType

	switch {
	case regexPost.MatchString(f.Url):
		matches := regexPost.FindStringSubmatch(f.Url)
		name := matches[1]
		id := matches[2]
		source = SourcePost{Id: id, name: name}
	case regexModel.MatchString(f.Url):
		matches := regexModel.FindStringSubmatch(f.Url)
		name := matches[1]
		source = SourceModel{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", f.Url)
	}

	f.Source = source
	return source, nil
}

func (f *Fapello) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Fapello)(nil)

// region - Private methods

func (f *Fapello) fetchMedia(
	source types.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var posts <-chan saktypes.Result[Post]

		switch s := source.(type) {
		case SourcePost:
			posts = f.fetchPost(s)
		case SourceModel:
			posts = f.fetchModel(s, limit)
		}

		for post := range posts {
			if post.Err != nil {
				out <- saktypes.Result[types.Media]{Err: post.Err}
				return
			}

			media := f.dataToMedia(post.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (f *Fapello) fetchPost(source SourcePost) <-chan saktypes.Result[Post] {
	result := make(chan saktypes.Result[Post])

	go func() {
		defer close(result)

		link := fmt.Sprintf("https://fapello.com/%s/%s", source.name, source.Id)
		post, err := getPost(link, source.name)

		if err != nil {
			result <- saktypes.Result[Post]{Err: err}
		} else {
			result <- saktypes.Result[Post]{Data: *post}
		}
	}()

	return result
}

func (f *Fapello) fetchModel(source SourceModel, limit int) <-chan saktypes.Result[Post] {
	result := make(chan saktypes.Result[Post])

	go func() {
		defer close(result)

		links, err := getLinks(source.name, limit)
		if err != nil {
			result <- saktypes.Result[Post]{Err: err}
			return
		}

		for _, link := range links {
			post, postErr := getPost(link, source.name)

			if postErr != nil {
				result <- saktypes.Result[Post]{Err: postErr}
			} else {
				result <- saktypes.Result[Post]{Data: *post}
			}
		}
	}()

	return result
}

func (f *Fapello) dataToMedia(post Post, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := f.DownloadHeaders()

	go func() {
		defer close(out)
		now := time.Date(1980, time.October, 6, 17, 7, 0, 0, time.UTC)

		media, err := types.NewMedia(post.Url, types.Fapello, map[string]interface{}{
			"id":      post.Id,
			"name":    post.Name,
			"source":  strings.ToLower(sourceName),
			"created": now.Add(time.Duration(post.Id*24) * time.Hour),
		}, headers)
		if err != nil {
			return
		}
		out <- media
	}()

	return out
}

// endregion
