package imaglr

import "time"

type Post struct {
	Id        string
	Author    string
	Type      string
	Media     string
	Thumbnail string
	Timestamp time.Time
}
