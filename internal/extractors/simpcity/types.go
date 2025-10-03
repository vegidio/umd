package simpcity

import "time"

type Post struct {
	Id          string
	Url         string
	Name        string
	Attachments []Attachment
	Published   time.Time
}

type Attachment struct {
	MediaUrl string
	ThumbUrl string
}
