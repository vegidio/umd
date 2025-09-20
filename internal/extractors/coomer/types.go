package coomer

import (
	"encoding/json"
	"fmt"

	"github.com/vegidio/go-sak/time"
)

type Response struct {
	Post   *Post  `json:"post"`
	Images []File `json:"previews"`
	Videos []File `json:"attachments"`
	Props  *Props `json:"props"`
}

type Profile struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Service   string `json:"service"`
	PostCount int    `json:"post_count"`
}

type Post struct {
	Id          string        `json:"id"`
	RevisionId  int           `json:"revision_id"`
	Service     string        `json:"service"`
	User        string        `json:"user"`
	Published   time.NotzTime `json:"published"`
	Attachments []File        `json:"attachments"`
}

type File struct {
	Server string `json:"server"`
	Path   string `json:"path"`
}

type Props struct {
	Revisions []Revision `json:"revisions"`
}

type Revision struct {
	Post Post `json:"-"`
}

type ResponseRevision struct {
	Post   *Post  `json:"post"`
	Images []File `json:"result_previews"`
	Videos []File `json:"result_attachments"`
	Props  *Props `json:"props"`
}

func (r *Revision) UnmarshalJSON(b []byte) error {
	var pair []json.RawMessage
	if err := json.Unmarshal(b, &pair); err != nil {
		return err
	}
	if len(pair) != 2 {
		return fmt.Errorf("revision must be [version, post], got %d items", len(pair))
	}
	if err := json.Unmarshal(pair[1], &r.Post); err != nil {
		return fmt.Errorf("post: %w", err)
	}
	return nil
}
