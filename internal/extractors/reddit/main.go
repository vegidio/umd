package reddit

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

const Host = "reddit.com"

type Reddit struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	headers          map[string]string
	external         types.External
}

func New(url string, metadata types.Metadata, headers map[string]string, external types.External) types.Extractor {
	switch {
	case utils.HasHost(url, Host):
		return &Reddit{Metadata: metadata, url: url, headers: headers, external: external}
	}

	return nil
}

func (r *Reddit) Type() types.ExtractorType {
	return types.Reddit
}

func (r *Reddit) SourceType() (types.SourceType, error) {
	regexSubmission := regexp.MustCompile(`/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/(?:u|user)/([^/\n?]+)`)
	regexSubreddit := regexp.MustCompile(`/r/([^/\n]+)`)

	var source types.SourceType
	var name string

	switch {
	case regexSubmission.MatchString(r.url):
		matches := regexSubmission.FindStringSubmatch(r.url)
		name = matches[1]
		id := matches[2]
		source = SourceSubmission{Id: id, name: name}

	case regexUser.MatchString(r.url):
		matches := regexUser.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceUser{name: name}

	case regexSubreddit.MatchString(r.url):
		matches := regexSubreddit.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceSubreddit{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.url)
	}

	r.source = source
	return source, nil
}

func (r *Reddit) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if r.responseMetadata == nil {
		r.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       r.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Reddit,
		Metadata:  r.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if r.source == nil {
			r.source, err = r.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := r.fetchMedia(r.source, extensions, deep)

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
var _ types.Extractor = (*Reddit)(nil)

// region - Private methods

func (r *Reddit) fetchMedia(
	source types.SourceType,
	extensions []string,
	deep bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

	go func() {
		defer close(out)
		var children <-chan saktypes.Result[ChildData]

		switch s := source.(type) {
		case SourceSubmission:
			children = getSubmission(s.Id)
		case SourceUser:
			children = getUserSubmissions(s.name)
		case SourceSubreddit:
			children = getSubredditSubmissions(s.name)
		}

		for child := range children {
			if child.Err != nil {
				out <- saktypes.Result[[]types.Media]{Err: child.Err}
				return
			}

			media := r.childToMedia(child.Data, source.Type(), source.Name())
			if deep {
				media = r.external.ExpandMedia(media, Host, &r.responseMetadata, 5)
			}

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

func (r *Reddit) childToMedia(child ChildData, sourceName string, name string) []types.Media {
	url := child.SecureMedia.RedditVideo.FallbackUrl
	if url == "" {
		url = child.Url
	}

	newMedia := types.NewMedia(url, types.Reddit, map[string]interface{}{
		"source":  strings.ToLower(sourceName),
		"name":    name,
		"created": child.Created.Time,
	})

	return []types.Media{newMedia}
}

// endregion
