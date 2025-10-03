package imaglr

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/samber/lo"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Imaglr struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "imaglr.com"):
		return &Imaglr{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (i *Imaglr) Type() types.ExtractorType {
	return types.Imaglr
}

func (i *Imaglr) SourceType() (types.SourceType, error) {
	regexPost := regexp.MustCompile(`/post/([^/\n?]+)`)

	var source types.SourceType
	var id string

	switch {
	case regexPost.MatchString(i.url):
		matches := regexPost.FindStringSubmatch(i.url)
		id = matches[1]
		source = SourcePost{name: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", i.url)
	}

	i.source = source
	return source, nil
}

func (i *Imaglr) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if i.responseMetadata == nil {
		i.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       i.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Imaglr,
		Metadata:  i.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if i.source == nil {
			i.source, err = i.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := i.fetchMedia(i.source, extensions, deep)

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
var _ types.Extractor = (*Imaglr)(nil)

// region - Private methods

func (i *Imaglr) fetchMedia(
	source types.SourceType,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

	go func() {
		defer close(out)

		posts := make([]Post, 0)
		var err error

		switch s := source.(type) {
		case SourcePost:
			posts, err = i.fetchPost(s)
		}

		if err != nil {
			out <- saktypes.Result[[]types.Media]{Err: err}
			return
		}

		media := postsToMedia(posts, source.Name())

		// Filter files with certain extensions
		if len(extensions) > 0 {
			media = lo.Filter(media, func(m types.Media, _ int) bool {
				return slices.Contains(extensions, m.Extension)
			})
		}

		out <- saktypes.Result[[]types.Media]{Data: media}
	}()

	return out
}

func (i *Imaglr) fetchPost(source SourcePost) ([]Post, error) {
	post, err := getPost(source.name)

	if err != nil {
		return make([]Post, 0), err
	}

	return []Post{*post}, nil
}

// endregion

// region - Private functions

func postsToMedia(posts []Post, sourceName string) []types.Media {
	return lo.Map(posts, func(post Post, _ int) types.Media {
		return types.NewMedia(post.Media, types.Imaglr, map[string]interface{}{
			"id":      post.Id,
			"name":    post.Author,
			"source":  strings.ToLower(sourceName),
			"created": post.Timestamp,
		})
	})
}

// endregion
