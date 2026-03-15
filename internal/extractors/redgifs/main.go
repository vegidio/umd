package redgifs

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Redgifs struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "redgifs.com"):
		r := &Redgifs{}
		r.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.RedGifs,
		}
		r.FetchMediaFn = r.fetchMedia
		r.SourceTypeFn = r.SourceType
		return r, nil
	}

	return nil, nil
}

var (
	regexVideo = regexp.MustCompile(`/(ifr|watch)/([^/\n?]+)`)
	regexUser  = regexp.MustCompile(`/users/([^/\n?]+)`)
)

func (r *Redgifs) SourceType() (types.SourceType, error) {

	var source types.SourceType
	var name string

	switch {
	case regexVideo.MatchString(r.Url):
		matches := regexVideo.FindStringSubmatch(r.Url)
		name = matches[2]
		source = SourceVideo{name: name}
	case regexUser.MatchString(r.Url):
		matches := regexUser.FindStringSubmatch(r.Url)
		name = matches[1]
		source = SourceUser{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.Url)
	}

	r.Source = source
	return source, nil
}

func (r *Redgifs) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Redgifs)(nil)

// region - Private methods

func (r *Redgifs) getNewOrSavedToken() (string, error) {
	token, exists := r.BaseExtractor.Metadata[types.RedGifs]["token"].(string)

	if !exists {
		log.Debug("Issuing new RedGifs token")

		auth, err := getToken()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to issue RedGifs token")

			return "", err
		}

		token = auth.Token

		if r.ResponseMetadata[types.RedGifs] == nil {
			r.ResponseMetadata[types.RedGifs] = make(map[string]interface{})
		}

		// Save the token to be reused in the future
		r.ResponseMetadata[types.RedGifs]["token"] = token
	} else {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("Reusing RedGifs token")
	}

	return token, nil
}

func (r *Redgifs) fetchMedia(
	source types.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var gifs <-chan saktypes.Result[[]Gif]

		token, err := r.getNewOrSavedToken()
		if err != nil {
			out <- saktypes.Result[types.Media]{Err: err}
			return
		}

		switch s := source.(type) {
		case SourceVideo:
			gifs = r.fetchGif(s, token)
		case SourceUser:
			gifs = r.fetchUser(s, token, limit)
		}

		for gif := range gifs {
			if gif.Err != nil {
				out <- saktypes.Result[types.Media]{Err: gif.Err}
				return
			}

			media := r.dataToMedia(gif.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (r *Redgifs) fetchGif(source SourceVideo, token string) <-chan saktypes.Result[[]Gif] {
	result := make(chan saktypes.Result[[]Gif])

	go func() {
		defer close(result)

		response, err := getGif(
			fmt.Sprintf("Bearer %s", token),
			fmt.Sprintf("https://www.redgifs.com/watch/%s", source.name),
			source.name,
		)

		if err != nil {
			result <- saktypes.Result[[]Gif]{Err: err}
			return
		}

		result <- saktypes.Result[[]Gif]{Data: []Gif{response.Gif}}
	}()

	return result
}

func (r *Redgifs) fetchUser(source SourceUser, token string, limit int) <-chan saktypes.Result[[]Gif] {
	result := make(chan saktypes.Result[[]Gif])

	go func() {
		defer close(result)

		bearer := fmt.Sprintf("Bearer %s", token)
		url := fmt.Sprintf("https://www.redgifs.com/users/%s", source.name)
		response, err := getUser(bearer, url, source.name, 1)

		if err != nil {
			result <- saktypes.Result[[]Gif]{Err: err}
			return
		}

		result <- saktypes.Result[[]Gif]{Data: response.Gifs}
		maxPages := math.Ceil(float64(limit) / 100)
		numPages := int(math.Min(float64(response.Pages), maxPages))

		for i := 2; i <= numPages; i++ {
			response, err = getUser(bearer, url, source.name, i)
			if err != nil {
				result <- saktypes.Result[[]Gif]{Err: err}
				return
			}

			result <- saktypes.Result[[]Gif]{Data: response.Gifs}
		}
	}()

	return result
}

func (r *Redgifs) dataToMedia(gifs []Gif, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := r.DownloadHeaders()

	go func() {
		defer close(out)

		for _, gif := range gifs {
			url := gif.Url.Hd
			if url == "" {
				url = gif.Url.Sd
			}

			media, err := types.NewMedia(url, types.RedGifs, map[string]interface{}{
				"name":    gif.Username,
				"source":  strings.ToLower(sourceName),
				"created": gif.Created.Time,
				"id":      gif.Id,
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
