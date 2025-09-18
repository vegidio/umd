package redgifs

import (
	"fmt"

	"github.com/vegidio/go-sak/fetch"
)

const BaseUrl = "https://api.redgifs.com/"

var f = fetch.New(nil, 0)

func getToken() (*Auth, error) {
	var auth *Auth
	url := BaseUrl + "v2/auth/temporary"
	headers := map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://www.redgifs.com",
		"Referer":      "https://www.redgifs.com/",
	}

	resp, err := f.GetResult(url, headers, &auth)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching authorization token: %s", resp.Status())
	}

	return auth, nil
}

func getGif(token string, videoUrl string, videoId string) (*GifResponse, error) {
	var response *GifResponse
	url := BaseUrl + fmt.Sprintf("v2/gifs/%s?views=yes", videoId)
	headers := map[string]string{
		"Authorization":  token,
		"X-CustomHeader": videoUrl,
	}

	resp, err := f.GetResult(url, headers, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching video ID '%s': %s", videoId, resp.Status())
	}

	return response, nil
}

func getUser(token string, userUrl string, userName string, page int) (*UserResponse, error) {
	var response *UserResponse
	url := BaseUrl + fmt.Sprintf("v2/users/%s/search?page=%d&count=100&order=latest&type=a&views=yes", userName, page)
	headers := map[string]string{
		"Authorization":  token,
		"X-CustomHeader": userUrl,
	}

	resp, err := f.GetResult(url, headers, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching username '%s': %s", userName, resp.Status())
	}

	return response, nil
}
