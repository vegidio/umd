package tests

import (
	"os"
	"testing"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
	"github.com/vegidio/umd/internal/types"
)

func TestSimpCity_QueryThread(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	log.SetLevel(log.DebugLevel)
	const NumberOfPosts = 471

	cookies, _ := fetch.GetFileCookies("/Users/vegidio/Desktop/cookies.txt")

	metadata := umd.Metadata{
		types.SimpCity: map[string]interface{}{
			"maxPages": 1,
			"cookie":   fetch.CookiesToHeader(cookies),
		},
	}

	extractor, _ := umd.New().
		WithMetadata(metadata).
		FindExtractor("https://simpcity.cr/threads/jessica-nigri.9946")

	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://simpcity.cr/attachments/1215x1688_7ec0a1bb6b3e911e892e54556be53825-jpg.2063/", resp.Media[0].Url)
	assert.Equal(t, "thread", resp.Media[0].Metadata["source"])
	assert.Equal(t, "Jessica Nigri", resp.Media[0].Metadata["name"])
}

func TestSimpCity_QueryThread_Page45(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	log.SetLevel(log.DebugLevel)
	const NumberOfPosts = 103

	cookies, _ := fetch.GetFileCookies("/Users/vegidio/Desktop/cookies.txt")

	metadata := umd.Metadata{
		types.SimpCity: map[string]interface{}{
			"startPage": 45,
			"maxPages":  1,
			"cookie":    fetch.CookiesToHeader(cookies),
		},
	}

	extractor, _ := umd.New().
		WithMetadata(metadata).
		FindExtractor("https://simpcity.cr/threads/jessica-nigri.9946")

	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	_, exists := lo.Find(resp.Media, func(m types.Media) bool {
		return m.Url == "https://simp6.selti-delivery.ru/images3/1000010158bb9dba913dcf1953.jpg"
	})

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "thread", resp.Media[0].Metadata["source"])
	assert.Equal(t, "Jessica Nigri", resp.Media[0].Metadata["name"])
	assert.True(t, exists)
}
