package bunkr

import "time"

type Response struct {
	Encrypted bool   `json:"encrypted"`
	Timestamp int    `json:"timestamp"`
	Url       string `json:"url"`
}

type Image struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Url       string `json:"auth_url"`
	Published time.Time
}

type Auth struct {
	Url string
}
