package reddit

import (
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
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, Host):
		r := &Reddit{}
		r.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Reddit,
		}
		r.FetchMediaFn = r.fetchMedia
		r.SourceTypeFn = r.SourceType
		return r, nil
	}

	return nil, nil
}

var (
	regexSubmission = regexp.MustCompile(`/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)`)
	regexUser       = regexp.MustCompile(`/(?:u|user)/([^/\n?]+)`)
	regexSubreddit  = regexp.MustCompile(`/r/([^/\n]+)`)
)

func (r *Reddit) SourceType() (types.SourceType, error) {

	var source types.SourceType
	var name string

	switch {
	case regexSubmission.MatchString(r.Url):
		matches := regexSubmission.FindStringSubmatch(r.Url)
		name = matches[1]
		id := matches[2]
		source = SourceSubmission{Id: id, name: name}

	case regexUser.MatchString(r.Url):
		matches := regexUser.FindStringSubmatch(r.Url)
		name = matches[1]
		source = SourceUser{name: name}

	case regexSubreddit.MatchString(r.Url):
		matches := regexSubreddit.FindStringSubmatch(r.Url)
		name = matches[1]
		source = SourceSubreddit{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.Url)
	}

	r.Source = source
	return source, nil
}

func (r *Reddit) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Reddit)(nil)

// region - Private methods

func (r *Reddit) getNewOrSavedToken() (string, error) {
	token, exists := r.BaseExtractor.Metadata[types.Reddit]["token"].(string)

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

		if r.ResponseMetadata[types.Reddit] == nil {
			r.ResponseMetadata[types.Reddit] = make(map[string]interface{})
		}

		// Save the token to be reused in the future
		r.ResponseMetadata[types.Reddit]["token"] = token
	} else {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("Reusing Reddit token")
	}

	return token, nil
}

func (r *Reddit) fetchMedia(
	source types.SourceType,
	_ int,
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
					return r.External.ExpandMedia(m, Host, &r.ResponseMetadata)
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

		media, err := types.NewMedia(url, types.Reddit, map[string]interface{}{
			"source":  strings.ToLower(sourceName),
			"name":    name,
			"created": child.Created.Time,
		}, headers)
		if err != nil {
			return
		}
		out <- media
	}()

	return out
}

// endregion
