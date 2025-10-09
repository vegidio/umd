package saint

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd/internal/utils"
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
		Published: utils.FakeTimestamp(id),
	}, nil
}
