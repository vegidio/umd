package coomer

import (
	"fmt"
	"regexp"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Coomer struct {
	types.BaseExtractor

	baseUrl   string
	regexPost *regexp.Regexp
	regexUser *regexp.Regexp
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "coomer.party", "coomer.st", "coomer.su"):
		services := "onlyfans|fansly|candfans"
		c := &Coomer{
			baseUrl:   "https://coomer.st",
			regexPost: regexp.MustCompile(`(` + services + `)/user/([^/]+)/post/([^/\n?]+)`),
			regexUser: regexp.MustCompile(`(` + services + `)/user/([^/\n?]+)`),
		}
		c.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Coomer,
		}
		c.FetchMediaFn = c.fetchMedia
		c.SourceTypeFn = c.SourceType
		return c, nil

	case utils.HasHost(url, "kemono.party", "kemono.su", "kemono.cr"):
		services := "patreon|fanbox|discord|fantia|afdian|boosty|gumroad|subscribestar|dlsite"
		c := &Coomer{
			baseUrl:   "https://kemono.cr",
			regexPost: regexp.MustCompile(`(` + services + `)/user/([^/]+)/post/([^/\n?]+)`),
			regexUser: regexp.MustCompile(`(` + services + `)/user/([^/\n?]+)`),
		}
		c.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Kemono,
		}
		c.FetchMediaFn = c.fetchMedia
		c.SourceTypeFn = c.SourceType
		return c, nil
	}

	return nil, nil
}

func (c *Coomer) SourceType() (types.SourceType, error) {
	var source types.SourceType
	var user string

	switch {
	case c.regexPost.MatchString(c.Url):
		matches := c.regexPost.FindStringSubmatch(c.Url)
		service := matches[1]
		user = matches[2]
		id := matches[3]
		source = SourcePost{Service: service, Id: id, name: user}

	case c.regexUser.MatchString(c.Url):
		matches := c.regexUser.FindStringSubmatch(c.Url)
		service := matches[1]
		user = matches[2]
		source = SourceUser{Service: service, name: user}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", c.Url)
	}

	c.Source = source
	return source, nil
}

func (c *Coomer) DownloadHeaders() map[string]string {
	headers := make(map[string]string)

	cookie, exists := c.Metadata[types.Coomer]["cookie"].(string)
	if exists {
		headers["Cookie"] = cookie
	}

	return headers
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Coomer)(nil)

// region - Private methods

func (c *Coomer) fetchMedia(
	source types.SourceType,
	_ int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	headers := make(map[string]string)
	cookie, exists := c.Metadata[types.Coomer]["cookie"].(string)
	if exists {
		headers["Cookie"] = cookie
	}

	profile, pErr := getProfile(c.baseUrl, source.(serviceSource).ServiceName(), source.Name(), headers)

	if pErr != nil {
		out <- saktypes.Result[types.Media]{Err: pErr}
		return out
	}

	go func() {
		defer close(out)
		var responses <-chan saktypes.Result[Response]

		switch s := source.(type) {
		case SourceUser:
			responses = getUser(c.baseUrl, *profile, headers)
		case SourcePost:
			responses = getPost(c.baseUrl, *profile, s.Id, headers)
		}

		for response := range responses {
			if response.Err != nil {
				out <- saktypes.Result[types.Media]{Err: response.Err}
				return
			}

			media := c.dataToMedia(response.Data, profile.Name)
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (c *Coomer) dataToMedia(response Response, name string) <-chan types.Media {
	out := make(chan types.Media)
	headers := c.DownloadHeaders()

	go func() {
		defer close(out)

		files := append(response.Images, response.Videos...)
		for _, file := range files {
			if file.Path != "" {
				url := file.Server + "/data" + file.Path
				media, err := types.NewMedia(url, c.ExtType, map[string]interface{}{
					"source":  response.Post.Service,
					"name":    name,
					"created": response.Post.Published.Time,
				}, headers)
				if err != nil {
					continue
				}
				out <- media
			}
		}
	}()

	return out
}

// endregion
