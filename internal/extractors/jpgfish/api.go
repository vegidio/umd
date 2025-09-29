package jpgfish

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/go-sak/fetch"
)

const BaseUrl = "https://jpg6.su/"

var f = fetch.New(nil, 0)

func getImage(id string) (*Image, error) {
	url := fmt.Sprintf("%simg/%s", BaseUrl, id)
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	title := doc.Find("a[data-text='image-title']").Text()
	encryptedUrl, _ := doc.Find("a.btn-download").Attr("href")
	author := doc.Find("a.user-link > strong").Text()

	timeObj := time.Now()
	timeStr, exists := doc.Find("p.description-meta > span").Attr("title")
	if exists {
		timeObj, _ = time.Parse("2006-01-02 15:04:05", timeStr)
	}

	return &Image{
		Id:        id,
		Title:     title,
		Url:       decryptUrl(encryptedUrl),
		Author:    author,
		Published: timeObj,
	}, nil
}

func decryptUrl(encryptedUrl string) string {
	// Hardcoded key
	key := []byte("seltilovessimpcity@simpcityhatesscrapers")

	// Base64 decode
	base64Decoded, err := base64.StdEncoding.DecodeString(encryptedUrl)
	if err != nil {
		return ""
	}

	// Convert from hex string to bytes
	hexDecoded, err := hex.DecodeString(string(base64Decoded))
	if err != nil {
		return ""
	}

	// XOR decryption
	keyLen := len(key)
	decrypted := make([]byte, len(hexDecoded))

	for i := 0; i < len(hexDecoded); i++ {
		decrypted[i] = hexDecoded[i] ^ key[i%keyLen]
	}

	// Convert back to string
	return string(decrypted)
}
