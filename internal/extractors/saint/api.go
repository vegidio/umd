package saint

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
)

const BaseUrl = "https://saint2.su/"

var f = fetch.New(nil, 0)

func getVideo(id string) (*Video, error) {
	url := fmt.Sprintf("%sembed/%s", BaseUrl, id)
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	videoUrl, _ := doc.Find("video > source").Attr("src")

	return &Video{
		Id:        id,
		Url:       videoUrl,
		Published: stringToTimeBase62(id),
	}, nil
}

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62ToBigInt(s string) (*big.Int, error) {
	n := big.NewInt(0)
	base := big.NewInt(62)

	for i := 0; i < len(s); i++ {
		c := s[i]
		idx := int64(strings.IndexByte(base62Alphabet, c))

		if idx < 0 {
			return nil, fmt.Errorf("invalid base62 char %q at pos %d", c, i)
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(idx))
	}
	return n, nil
}

func stringToTimeBase62(s string) time.Time {
	start := time.Date(1980, time.October, 6, 17, 7, 0, 0, time.UTC)
	end := time.Date(2035, 12, 31, 23, 59, 59, 0, time.UTC)

	if !end.After(start) {
		return time.Now()
	}
	n, err := base62ToBigInt(s)
	if err != nil {
		return time.Now()
	}

	// We map to seconds resolution within [start, end].
	rangeSeconds := end.Unix() - start.Unix() + 1 // +1 to include 'end' second
	if rangeSeconds <= 0 {
		return time.Now()
	}

	r := new(big.Int).Mod(n, big.NewInt(rangeSeconds)).Int64()
	return start.Add(time.Duration(r) * time.Second)
}
