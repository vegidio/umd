package erome

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd/internal/utils"
)

const BaseUrl = "https://www.erome.com/a/"

var f = fetch.New(nil, 0)

func getAlbum(id string) (*Album, error) {
	links := make([]string, 0)

	url := BaseUrl + id
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	name := doc.Find("h1.album-title-page").Text()
	user := doc.Find("a#user_name").Text()
	created := utils.FakeTimestamp(id)

	doc.Find("div.img > img").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("data-src"); exists {
			links = append(links, link)
		}
	})

	doc.Find("div.video-lg > video.video-js > source").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("src"); exists {
			links = append(links, link)
		}
	})

	return &Album{
		Id:      id,
		Name:    name,
		User:    user,
		Url:     url,
		Created: created,
		Links:   links,
	}, nil
}
