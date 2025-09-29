package jpgfish

import "time"

type Image struct {
	Id        string
	Title     string
	Url       string
	Author    string
	Published time.Time
}
