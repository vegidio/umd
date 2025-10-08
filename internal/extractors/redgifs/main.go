package redgifs

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Redgifs struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "redgifs.com"):
		return &Redgifs{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (r *Redgifs) Type() types.ExtractorType {
	return types.RedGifs
}

func (r *Redgifs) SourceType() (types.SourceType, error) {
	regexVideo := regexp.MustCompile(`/(ifr|watch)/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/users/([^/\n?]+)`)

	var source types.SourceType
	var name string

	switch {
	case regexVideo.MatchString(r.url):
		matches := regexVideo.FindStringSubmatch(r.url)
		name = matches[2]
		source = SourceVideo{name: name}
	case regexUser.MatchString(r.url):
		matches := regexUser.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceUser{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.url)
	}

	r.source = source
	return source, nil
}

func (r *Redgifs) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if r.responseMetadata == nil {
		r.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       r.url,
		Media:     make([]types.Media, 0),
		Extractor: types.RedGifs,
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

		mediaCh := r.fetchMedia(r.source, limit, extensions, deep)

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
var _ types.Extractor = (*Redgifs)(nil)

// region - Private methods

func (r *Redgifs) getNewOrSavedToken() (string, error) {
	token, exists := r.Metadata[types.RedGifs]["token"].(string)

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

		if r.responseMetadata[types.RedGifs] == nil {
			r.responseMetadata[types.RedGifs] = make(map[string]interface{})
		}

		// Save the token to be reused in the future
		r.responseMetadata[types.RedGifs]["token"] = token
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

			media := videosToMedia(gif.Data, source.Type())

			for m := range media {
				// Filter files with certain extensions
				if len(extensions) > 0 {
					if slices.Contains(extensions, m.Extension) {
						out <- saktypes.Result[types.Media]{Data: m}
						continue
					}
				}

				out <- saktypes.Result[types.Media]{Data: m}
			}
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

// endregion

// region - Private functions

func videosToMedia(gifs []Gif, sourceName string) <-chan types.Media {
	out := make(chan types.Media)

	go func() {
		defer close(out)

		for _, gif := range gifs {
			url := gif.Url.Hd
			if url == "" {
				url = gif.Url.Sd
			}

			out <- types.NewMedia(url, types.RedGifs, map[string]interface{}{
				"name":    gif.Username,
				"source":  strings.ToLower(sourceName),
				"created": gif.Created.Time,
				"id":      gif.Id,
			})
		}
	}()

	return out
}

// endregion
