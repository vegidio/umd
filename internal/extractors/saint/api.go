package saint

import (
	"fmt"

	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd/internal/utils"
)

const BaseUrl = "https://turbo.cr/"

var f = fetch.New(nil, 0)

func getVideo(id string) (*Video, error) {
	var response *Response
	url := fmt.Sprintf("%sapi/sign?v=%s", BaseUrl, id)
	resp, err := f.GetResult(url, map[string]string{
		"Referer": url,
	}, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, resp.Error().(error)
	}

	return &Video{
		Id:        id,
		Url:       response.Url,
		Published: utils.FakeTimestamp(id),
	}, nil
}
