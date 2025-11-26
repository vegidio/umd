package simpcity

import "time"

type Post struct {
	Id          string
	Url         string
	Title       string
	Attachments []Attachment
	Published   time.Time
}

type Attachment struct {
	MediaUrl string
	ThumbUrl string
}
