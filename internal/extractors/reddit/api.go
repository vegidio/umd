package reddit

import (
	"encoding/json"
	"fmt"

	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/go-sak/types"
)

const BaseUrl = "https://www.reddit.com/"

var f = fetch.New(nil, 10)

// getSubmission fetches and processes submission data for a given Reddit post ID.
//
// Example: https://www.reddit.com/comments/1bxsmnr.json?raw_json=1, where <1bxsmnr> is the ID.
//
// # Parameters:
//   - id: string - The unique identifier of the Reddit post to fetch
//
// # Returns:
//   - <-chan model.Result[ChildData] - A receive-only channel that streams Reddit post data or errors
func getSubmission(id string) <-chan types.Result[ChildData] {
	out := make(chan types.Result[ChildData])

	go func() {
		defer close(out)

		submissions := make([]Submission, 0)
		url := fmt.Sprintf(BaseUrl+"comments/%s.json?raw_json=1", id)
		resp, err := f.GetResult(url, nil, &submissions)

		if err != nil {
			out <- types.Result[ChildData]{Err: err}
			return
		} else if resp.IsError() {
			out <- types.Result[ChildData]{Err: fmt.Errorf("error fetching post id '%s' submissions: %s", id, resp.Status())}
			return
		}

		submission := submissions[0]
		for _, child := range submission.Data.Children {
			if child.Data.IsGallery {
				children := getGalleryData(child.Data)

				for _, gallery := range children {
					out <- types.Result[ChildData]{Data: gallery}
				}
			} else {
				out <- types.Result[ChildData]{Data: child.Data}
			}
		}
	}()

	return out
}

// getUserSubmissions retrieves a stream of user submissions as a channel of types.Result[ChildData]. The submissions
// are fetched using the specified user's name.
//
// Example: https://www.reddit.com/user/atomicbrunette18/submitted.json?sort=new&raw_json=1&after=&limit=100, where
// <atomicbrunette18> is the username.
//
// # Parameters:
//   - user: string - The username whose submissions to fetch
//
// # Returns:
//   - <-chan model.Result[ChildData] - A receive-only channel that streams submission data or errors
func getUserSubmissions(user string) <-chan types.Result[ChildData] {
	urlFmt := BaseUrl + "user/%s/submitted.json?sort=new&raw_json=1&after=%s&limit=%d"
	return streamSubmissions(urlFmt, user)
}

// getSubredditSubmissions retrieves a stream of subreddit submissions as a channel of types.Result[ChildData]. The
// submissions are fetched using the specified subreddit's name.
//
// Example: https://www.reddit.com/r/nsfw/hot.json?raw_json=1&after=&limit=100, where <nsfw> is the subreddit name.
//
// # Parameters:
//   - subreddit: string - The subreddit whose submissions are to fetch.
//
// # Returns:
//   - <-chan types.Result[ChildData] - A receive-only channel that streams submission data or errors.
func getSubredditSubmissions(subreddit string) <-chan types.Result[ChildData] {
	urlFmt := BaseUrl + "r/%s/hot.json?raw_json=1&after=%s&limit=%d"
	return streamSubmissions(urlFmt, subreddit)
}

func streamSubmissions(urlFmt string, what string) <-chan types.Result[ChildData] {
	out := make(chan types.Result[ChildData])

	go func() {
		defer close(out)
		after := ""

		for {
			var submission *Submission
			url := fmt.Sprintf(urlFmt, what, after, 100)
			resp, err := f.GetResult(url, nil, &submission)

			if err != nil {
				out <- types.Result[ChildData]{Err: err}
				return
			} else if resp.IsError() {
				out <- types.Result[ChildData]{Err: fmt.Errorf("error fetching %s submissions: %s", what, resp.Status())}
				return
			}

			for _, child := range submission.Data.Children {
				if child.Data.IsGallery {
					for _, galleryItem := range getGalleryData(child.Data) {
						out <- types.Result[ChildData]{Data: galleryItem}
					}
				} else {
					out <- types.Result[ChildData]{Data: child.Data}
				}
			}

			after = submission.Data.After
			if after == "" {
				return
			}
		}
	}()
	return out
}

func getGalleryData(child ChildData) []ChildData {
	children := make([]ChildData, 0)

	for _, value := range child.MediaMetadata {
		var metadata MediaMetadata
		jsonData, _ := json.Marshal(value)
		json.Unmarshal(jsonData, &metadata)

		if metadata.Status == "valid" {
			url := metadata.S.Image
			if url == "" {
				url = metadata.S.Gif
			}

			newChild := ChildData{
				Author:  child.Author,
				Url:     url,
				Created: child.Created,
			}

			children = append(children, newChild)
		}
	}

	return children
}
