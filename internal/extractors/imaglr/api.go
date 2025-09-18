package imaglr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
)

const BaseUrl = "https://imaglr.com/"

var f = fetch.New(nil, 0)

func getPost(id string) (*Post, error) {
	url := BaseUrl + fmt.Sprintf("post/%s", id)
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var result map[string]any
	jsonString, _ := doc.Find("div#app").Attr("data-page")
	err = json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}

	data := result["props"].(map[string]any)["post"].(map[string]any)["data"].(map[string]any)
	media := data["media"].([]any)[0].(map[string]any)

	author := data["user"].(map[string]any)["name"].(string)
	mediaType := media["type"].(string)
	mediaUrl := media["media_url"].(string)
	thumbUrl := media["thumb_url"].(string)

	timestamp := result["props"].(map[string]interface{})["post"].(map[string]interface{})["data"].(map[string]interface{})["created_at_timestamp"].(float64)
	createdAt := time.Unix(int64(timestamp), 0)

	post := &Post{
		Id:        id,
		Author:    author,
		Type:      mediaType,
		Media:     mediaUrl,
		Thumbnail: thumbUrl,
		Timestamp: createdAt,
	}

	return post, nil
}
