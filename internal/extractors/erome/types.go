package erome

import "time"

type Album struct {
	Id      string
	Name    string
	User    string
	Url     string
	Created time.Time
	Links   []string
}
