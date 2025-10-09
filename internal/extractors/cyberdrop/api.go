package cyberdrop

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd/internal/utils"
)

const BaseUrl = "https://cyberdrop.me/"

var f = fetch.New(nil, 0)

func getImage(id string) (*Image, error) {
	var image *Image
	url := fmt.Sprintf("https://api.cyberdrop.me/api/file/info/%s", id)
	resp, err := f.GetResult(url, nil, &image)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, resp.Error().(error)
	}

	var auth *Auth
	resp, err = f.GetResult(image.Url, nil, &auth)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, resp.Error().(error)
	}

	image.Url = auth.Url
	image.Published = utils.FakeTimestamp(image.Id)

	return image, nil
}

func getAlbum(id string) ([]string, error) {
	ids := make([]string, 0)

	url := fmt.Sprintf("%sa/%s", BaseUrl, id)
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	doc.Find("a.image").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		imageId := strings.Split(href, "/")[2]
		ids = append(ids, imageId)
	})

	return ids, nil
}
