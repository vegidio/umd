package fapello

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
)

const BaseUrl = "https://fapello.com/"

var f = fetch.New(nil, 0)

func getLinks(name string, limit int) ([]string, error) {
	links := make([]string, 0)
	numPages := 1

	url := BaseUrl + name
	html, err := f.GetText(url)
	if err != nil {
		return links, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return links, err
	}

	showMore := doc.Find("div#showmore")
	if showMore.Length() > 0 {
		maxPages := math.Ceil(float64(limit) / 32)

		pages, _ := doc.Find("div#showmore").Attr("data-max")
		pagesF, errF := strconv.ParseFloat(pages, 64)
		if errF != nil {
			return links, errF
		}

		numPages = int(math.Min(pagesF, maxPages))
	}

	for i := 1; i <= numPages; i++ {
		pageUrl := fmt.Sprintf("%s/ajax/model/%s/page-%d/", BaseUrl, name, i)
		html, err = f.GetText(pageUrl)
		if err != nil {
			return links, err
		}

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return links, err
		}

		doc.Find("img.object-cover").Each(func(i int, s *goquery.Selection) {
			parentLink, _ := s.ParentsFiltered("a").Attr("href")
			links = append(links, parentLink)
		})
	}

	return links, nil
}

func getPost(url string, name string) (*Post, error) {
	mediaUrl := ""

	matches := regexp.MustCompile(`/(\d+)/?$`).FindStringSubmatch(url)
	id, _ := strconv.Atoi(matches[1])

	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	videoTag := doc.Find("video.uk-align-center")

	if videoTag.Length() > 0 {
		mediaUrl, _ = videoTag.Find("source").Attr("src")
	} else {
		mediaUrl, _ = doc.Find("div.flex.justify-between.items-center > a").Attr("href")
	}

	return &Post{
		Id:   id,
		Name: name,
		Url:  mediaUrl,
	}, nil
}
