package coomer

import "github.com/vegidio/go-sak/time"

type Response struct {
	Post   *Post  `json:"post"`
	Images []File `json:"previews"`
	Videos []File `json:"attachments"`
}

type Profile struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Service   string `json:"service"`
	PostCount int    `json:"post_count"`
}

type Post struct {
	Id        string        `json:"id"`
	Service   string        `json:"service"`
	User      string        `json:"user"`
	Published time.NotzTime `json:"published"`
}

type File struct {
	Server string `json:"server"`
	Path   string `json:"path"`
}
