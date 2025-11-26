package main

import (
	"fmt"
	"os"
	"shared"

	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

func main() {
	cookies, err := shared.GetCookies("manual", "/Users/vegidio/Desktop/cookies.txt")
	if err != nil {
		panic(err)
	}

	metadata := umd.Metadata{
		umd.SimpCity: map[string]interface{}{
			"cookie": fetch.CookiesToHeader(cookies),
		},
	}

	u := umd.New().WithMetadata(metadata)

	extractor, err := u.FindExtractor("https://simpcity.cr/threads/sylvia-yasmina-sylviayasmina.239680/")
	if err != nil {
		panic(err)
	}

	resp, _ := extractor.QueryMedia(10, nil, true)
	resp.Track(func(queried, total int) {
		fmt.Println("queried", queried)
	})

	os.MkdirAll("/Users/vegidio/Desktop/test", 0755)
	result := shared.DownloadAll(resp.Media, "/Users/vegidio/Desktop/test", 5)

	for file := range result {
		fmt.Println("Downloading", file.Request.Url)

		err = file.Error()
		if err == nil {
			fmt.Println("Done")
		} else {
			fmt.Println("Download", file.Request.Url, "failed")
		}
	}
}
