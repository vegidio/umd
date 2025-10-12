package bunkr

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd/internal/utils"
)

const BaseUrl = "https://bunkr.cr/"

var f = fetch.New(nil, 0)

func getImage(slug string) (*Image, error) {
	var response *Response
	url := "https://bunkr.cr/api/vs"
	headers := map[string]string{
		"Referer": "https://bunkr.cr/",
	}
	body := map[string]string{
		"slug": slug,
	}

	resp, err := f.PostResult(url, headers, body, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, resp.Error().(error)
	}

	mediaUrl := response.Url

	// The URL needs to be decrypted
	if response.Encrypted {
		mediaUrl, err = decryptUrl(mediaUrl, response.Timestamp)
		if err != nil {
			return nil, err
		}
	}

	fileRegex := regexp.MustCompile(`(-[^.]+)`)
	fileName := fileRegex.ReplaceAllString(filepath.Base(mediaUrl), "")

	return &Image{
		Slug:      slug,
		Name:      fileName,
		Url:       mediaUrl,
		Published: utils.FakeTimestamp(slug),
	}, nil
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

	doc.Find("div.theItem a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		imageId := strings.Split(href, "/")[2]
		ids = append(ids, imageId)
	})

	return ids, nil
}

func decryptUrl(encryptedURL string, timestamp int) (string, error) {
	// Generate a time-based key (hour-based timestamp)
	keyString := fmt.Sprintf("SECRET_KEY_%d", timestamp/3600)
	key := []byte(keyString)

	// Decode base64 encrypted data
	data, err := base64.StdEncoding.DecodeString(encryptedURL)
	if err != nil {
		return "", err
	}

	// XOR decryption with a repeating key
	keyLen := len(key)
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%keyLen]
	}

	decryptedURL := string(result)

	// Basic URL validation
	if len(decryptedURL) < 7 || (decryptedURL[:7] != "http://" && (len(decryptedURL) < 8 || decryptedURL[:8] != "https://")) {
		return "", fmt.Errorf("decrypted URL appears to be invalid: %s", decryptedURL)
	}

	return decryptedURL, nil
}
