package simpcity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/go-sak/types"
)

const BaseUrl = "https://simpcity.cr"

func getThread(id string, startPage, maxPages int, headers map[string]string) <-chan types.Result[Post] {
	out := make(chan types.Result[Post])

	go func() {
		defer close(out)

		f := fetch.New(headers, 0)
		url := fmt.Sprintf("%s/threads/%s", BaseUrl, id)
		html, err := f.GetText(url)
		if err != nil {
			out <- types.Result[Post]{Err: err}
			return
		}

		// Get the number of pages
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			out <- types.Result[Post]{Err: err}
			return
		}

		name := doc.Find("h1.p-title-value").Contents().Not("a, span").Text()
		name = strings.TrimSpace(name)

		pagesStr := doc.Find("li.pageNav-page > a").Last().Text()
		pagesNum, err := strconv.Atoi(pagesStr)
		if err != nil {
			out <- types.Result[Post]{Err: err}
			return
		}

		lastPage := startPage + maxPages - 1
		if lastPage > pagesNum || maxPages == 0 {
			lastPage = pagesNum
		}

		// Iterate through all pages
		for i := startPage; i <= lastPage; i++ {
			pageUrl := fmt.Sprintf("%s/page-%d", url, i)

			log.WithFields(log.Fields{
				"url": pageUrl,
			}).Debug("Parsing page")

			pageHtml, pErr := f.GetText(pageUrl)
			if pErr != nil {
				out <- types.Result[Post]{Err: pErr}
				return
			}

			pDoc, pErr := goquery.NewDocumentFromReader(strings.NewReader(pageHtml))
			if pErr != nil {
				out <- types.Result[Post]{Err: pErr}
				return
			}

			pDoc.Find("div.message-cell--main").Each(func(i int, q *goquery.Selection) {
				post, postErr := parsePost(id, name, q)

				if postErr != nil {
					out <- types.Result[Post]{Err: postErr}
				} else {
					out <- types.Result[Post]{Data: *post}
				}
			})
		}
	}()

	return out
}

func parsePost(id, name string, query *goquery.Selection) (*Post, error) {
	attachments := make([]Attachment, 0)
	timeS := query.Find("time.u-dt")
	published := time.Now()

	postUrl, exists := timeS.ParentFiltered("a").Attr("href")
	if exists {
		postUrl = BaseUrl + postUrl
	}

	timeStr, exists := timeS.Attr("data-timestamp")
	if exists {
		timestamp, err := strconv.ParseInt(timeStr, 10, 64)
		if err == nil {
			published = time.Unix(timestamp, 0)
		}
	}

	attachments = append(attachments, getAttachFilePreviewImage(query)...)
	attachments = append(attachments, getAttachFilePreviewVideo(query)...)
	attachments = append(attachments, getAttachBbImageWrapper(query)...)
	attachments = append(attachments, getAttachBbVideoWrapper(query)...)
	attachments = append(attachments, getAttachJsLbImage(query)...)
	attachments = append(attachments, getAttachSaintIFrame(query)...)
	attachments = append(attachments, getAttachLinkExternal(query)...)

	return &Post{
		Id:          id,
		Url:         postUrl,
		Name:        name,
		Attachments: attachments,
		Published:   published,
	}, nil
}

func trimAfterLastSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	lastSlash := strings.LastIndex(url, "/")
	if lastSlash == -1 {
		return url
	}

	return url[:lastSlash+1]
}

func getAttachFilePreviewImage(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("a.file-preview > img").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, exists := q.ParentFiltered("a").Attr("href")
		if exists && !strings.HasPrefix(mediaUrl, "https") {
			mediaUrl = BaseUrl + mediaUrl
		}

		thumbUrl, exists := q.Attr("alt")
		if exists && !strings.HasPrefix(thumbUrl, "https") {
			thumbUrl = trimAfterLastSlash(mediaUrl) + thumbUrl
		}

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: thumbUrl,
		})
	})

	return attachments
}

func getAttachFilePreviewVideo(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("a.file-preview > video").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, _ := q.ChildrenFiltered("source").Attr("src")

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: mediaUrl,
		})
	})

	return attachments
}

func getAttachBbImageWrapper(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("div.bbImageWrapper").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, exists := q.Attr("data-src")
		if exists && !strings.HasPrefix(mediaUrl, "https") {
			mediaUrl = BaseUrl + mediaUrl
		}

		thumbUrl, exists := q.Attr("title")
		if exists && !strings.HasPrefix(thumbUrl, "https") {
			thumbUrl = trimAfterLastSlash(mediaUrl) + thumbUrl
		}

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: thumbUrl,
		})
	})

	return attachments
}

func getAttachBbVideoWrapper(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("div.bbMediaWrapper source").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, _ := q.Attr("src")

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: mediaUrl,
		})
	})

	return attachments
}

func getAttachJsLbImage(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("a.js-lbImage").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, exists := q.Attr("href")
		if exists && !strings.HasPrefix(mediaUrl, "https") {
			mediaUrl = BaseUrl + mediaUrl
		}

		thumbUrl, exists := q.ChildrenFiltered("img").Attr("alt")
		if exists && !strings.HasPrefix(thumbUrl, "https") {
			thumbUrl = trimAfterLastSlash(mediaUrl) + thumbUrl
		}

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: thumbUrl,
		})
	})

	return attachments
}

func getAttachSaintIFrame(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("iframe.saint-iframe").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, _ := q.Attr("src")

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: mediaUrl,
		})
	})

	return attachments
}

func getAttachLinkExternal(query *goquery.Selection) []Attachment {
	attachments := make([]Attachment, 0)

	query.Find("a.link--external").Each(func(_ int, q *goquery.Selection) {
		mediaUrl, _ := q.Attr("href")

		attachments = append(attachments, Attachment{
			MediaUrl: mediaUrl,
			ThumbUrl: mediaUrl,
		})
	})

	return attachments
}
