package saint

import "time"

type Response struct {
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

type Video struct {
	Id        string
	Url       string
	Published time.Time
}
