package redgifs

import "github.com/vegidio/go-sak/time"

type Auth struct {
	Token   string `json:"token"`
	Session string `json:"session"`
}

type Url struct {
	Poster     string `json:"poster"`
	Thumbnail  string `json:"thumbnail"`
	Vthumbnail string `json:"vthumbnail"`
	Hd         string `json:"hd"`
	Sd         string `json:"sd"`
}

type Gif struct {
	Id       string         `json:"id"`
	Username string         `json:"userName"`
	Duration float64        `json:"duration"`
	Url      Url            `json:"urls"`
	Created  time.EpochTime `json:"createDate"`
}

type GifResponse struct {
	Gif Gif `json:"gif"`
}

type UserResponse struct {
	Page  int   `json:"page"`
	Pages int   `json:"pages"`
	Gifs  []Gif `json:"gifs"`
}
