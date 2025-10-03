package tests

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
	"github.com/vegidio/umd/internal/types"
)

func TestSimpCity_QueryThread(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	const NumberOfPosts = 183

	metadata := umd.Metadata{
		types.SimpCity: map[string]interface{}{
			"maxPages": 1,
		},
	}

	cookies, _ := fetch.GetFileCookies("/Users/vegidio/Desktop/cookies.txt")
	headers := map[string]string{
		"Cookie": fetch.CookiesToHeader(cookies),
	}

	extractor, _ := umd.New().
		WithMetadata(metadata).
		WithHeaders(headers).
		FindExtractor("https://simpcity.cr/threads/jessica-nigri.9946")

	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://simpcity.cr/attachments/1215x1688_7ec0a1bb6b3e911e892e54556be53825-jpg.2063/", resp.Media[0].Url)
	assert.Equal(t, "thread", resp.Media[0].Metadata["source"])
	assert.Equal(t, "Jessica Nigri", resp.Media[0].Metadata["name"])
}
