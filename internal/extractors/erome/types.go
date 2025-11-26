package erome

import "time"

type Album struct {
	Id      string
	Title   string
	User    string
	Created time.Time
	Links   []string
}
