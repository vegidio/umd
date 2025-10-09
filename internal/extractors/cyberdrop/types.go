package cyberdrop

import "time"

type Image struct {
	Id        string `json:"slug"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Url       string `json:"auth_url"`
	Published time.Time
}

type Auth struct {
	Url string
}
