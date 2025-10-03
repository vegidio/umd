package fapello

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Fapello struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	headers          map[string]string
	external         types.External
}

func New(url string, metadata types.Metadata, headers map[string]string, external types.External) types.Extractor {
	switch {
	case utils.HasHost(url, "fapello.com"):
		return &Fapello{Metadata: metadata, url: url, headers: headers, external: external}
	}

	return nil
}

func (f *Fapello) Type() types.ExtractorType {
	return types.Fapello
}

func (f *Fapello) SourceType() (types.SourceType, error) {
	regexPost := regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/(\d+)`)
	regexModel := regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/?`)

	var source types.SourceType

	switch {
	case regexPost.MatchString(f.url):
		matches := regexPost.FindStringSubmatch(f.url)
		name := matches[1]
		id := matches[2]
		source = SourcePost{Id: id, name: name}
	case regexModel.MatchString(f.url):
		matches := regexModel.FindStringSubmatch(f.url)
		name := matches[1]
		source = SourceModel{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", f.url)
	}

	f.source = source
	return source, nil
}

func (f *Fapello) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if f.responseMetadata == nil {
		f.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       f.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Fapello,
		Metadata:  f.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if f.source == nil {
			f.source, err = f.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := f.fetchMedia(f.source, limit, extensions, deep)

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

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Fapello)(nil)

// region - Private methods

func (f *Fapello) fetchMedia(
	source types.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

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
				out <- saktypes.Result[[]types.Media]{Err: post.Err}
				return
			}

			media := postsToMedia(post.Data, source.Type())

			// Filter files with certain extensions
			if len(extensions) > 0 {
				media = lo.Filter(media, func(m types.Media, _ int) bool {
					return slices.Contains(extensions, m.Extension)
				})
			}

			out <- saktypes.Result[[]types.Media]{Data: media}
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

// endregion

// region - Private functions

func postsToMedia(post Post, sourceName string) []types.Media {
	now := time.Date(1980, time.October, 6, 17, 7, 0, 0, time.UTC)

	return []types.Media{types.NewMedia(post.Url, types.Fapello, map[string]interface{}{
		"id":      post.Id,
		"name":    post.Name,
		"source":  strings.ToLower(sourceName),
		"created": now.Add(time.Duration(post.Id*24) * time.Hour),
	})}
}

// endregion
