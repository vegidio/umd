package reddit

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vegidio/go-sak/async"
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
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, Host):
		return &Reddit{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
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

func (r *Reddit) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Reddit)(nil)

// region - Private methods

func (r *Reddit) getNewOrSavedToken() (string, error) {
	token, exists := r.Metadata[types.Reddit]["token"].(string)

	if !exists {
		log.Debug("Issuing new Reddit token")

		auth, err := getToken()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to issue Reddit token")

			return "", err
		}

		token = auth.Token

		if r.responseMetadata[types.Reddit] == nil {
			r.responseMetadata[types.Reddit] = make(map[string]interface{})
		}

		// Save the token to be reused in the future
		r.responseMetadata[types.Reddit]["token"] = token
	} else {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("Reusing Reddit token")
	}

	return token, nil
}

func (r *Reddit) fetchMedia(
	source types.SourceType,
	extensions []string,
	deep bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var children <-chan saktypes.Result[ChildData]

		token, err := r.getNewOrSavedToken()
		if err != nil {
			out <- saktypes.Result[types.Media]{Err: err}
			return
		}

		switch s := source.(type) {
		case SourceSubmission:
			children = getSubmission(s.Id, token)
		case SourceUser:
			children = getUserSubmissions(s.name, token)
		case SourceSubreddit:
			children = getSubredditSubmissions(s.name, token)
		}

		for child := range children {
			if child.Err != nil {
				out <- saktypes.Result[types.Media]{Err: child.Err}
				return
			}

			media := r.dataToMedia(child.Data, source.Type(), source.Name())
			if deep {
				media = async.ConcurrentChannel(media, 5, func(m types.Media) types.Media {
					return r.external.ExpandMedia(m, Host, &r.responseMetadata)
				})
			}

			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (r *Reddit) dataToMedia(child ChildData, sourceName string, name string) <-chan types.Media {
	out := make(chan types.Media)
	headers := r.DownloadHeaders()

	go func() {
		defer close(out)

		url := child.SecureMedia.RedditVideo.FallbackUrl
		if url == "" {
			url = child.Url
		}

		out <- types.NewMedia(url, types.Reddit, map[string]interface{}{
			"source":  strings.ToLower(sourceName),
			"name":    name,
			"created": child.Created.Time,
		}, headers)
	}()

	return out
}

// endregion
