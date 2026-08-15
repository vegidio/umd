package coomer

import (
	"context"
	"fmt"
	"maps"

	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/utils"
)

var f = fetch.New(nil, 10, false)
var cssHeaders = map[string]string{"Accept": "text/css"}

func getProfile(ctx context.Context, baseUrl string, service string, user string, headers map[string]string) (*Profile, error) {
	maps.Copy(headers, cssHeaders)

	var profile *Profile
	url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/profile", service, user)
	resp, err := f.GetResult(ctx, url, headers, &profile)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching user '%s' posts: %s", user, resp.Status())
	}

	return profile, nil
}

func getUser(ctx context.Context, baseUrl string, profile Profile, headers map[string]string) <-chan types.Result[Response] {
	out := make(chan types.Result[Response])

	maps.Copy(headers, cssHeaders)

	go func() {
		defer close(out)

		for offset := 0; offset <= profile.PostCount; offset += 50 {
			var posts []Post
			url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/posts?o=%d", profile.Service, profile.Id, offset)
			resp, err := f.GetResult(ctx, url, headers, &posts)

			if err == nil && resp.IsError() {
				err = fmt.Errorf("error fetching user '%s' posts: %s", profile.Name, resp.Status())
			}

			if err != nil && !utils.Send(ctx, out, types.Result[Response]{Err: err}) {
				return
			}

			if len(posts) == 0 {
				break
			}

			for _, post := range posts {
				result := <-getPost(ctx, baseUrl, profile, post.Id, headers)
				if result.Err != nil {
					if !utils.Send(ctx, out, types.Result[Response]{Err: result.Err}) {
						return
					}
					continue
				}

				if !utils.Send(ctx, out, types.Result[Response]{Data: result.Data}) {
					return
				}
			}

			offset += 50
		}
	}()

	return out
}

func getPost(ctx context.Context, baseUrl string, profile Profile, postId string, headers map[string]string) <-chan types.Result[Response] {
	out := make(chan types.Result[Response])

	maps.Copy(headers, cssHeaders)

	go func() {
		defer close(out)

		var response Response
		url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/post/%s", profile.Service, profile.Id, postId)
		resp, err := f.GetResult(ctx, url, headers, &response)

		if err != nil {
			utils.Send(ctx, out, types.Result[Response]{Err: err})
			return
		} else if resp.IsError() {
			utils.Send(ctx, out, types.Result[Response]{Err: fmt.Errorf("error fetching user '%s' posts: %s",
				profile.Name, resp.Status())})
			return
		}

		// Check if the post has a revision with more attachments than the current post
		// If so, fetch the revision and return it instead of the current post
		biggestRevision := lo.MaxBy(response.Props.Revisions, func(a, b Revision) bool {
			return len(a.Post.Attachments) > len(b.Post.Attachments)
		})

		if biggestRevision.Post.RevisionId > 0 && (len(response.Images)+len(response.Videos)) < len(biggestRevision.Post.Attachments) {
			utils.Send(ctx, out, getRevision(ctx, baseUrl, profile, postId, biggestRevision.Post.RevisionId, headers))
			return
		}

		utils.Send(ctx, out, types.Result[Response]{Data: response})
	}()

	return out
}

func getRevision(ctx context.Context, baseUrl string, profile Profile, postId string, revisionId int, headers map[string]string) types.Result[Response] {
	var response ResponseRevision
	url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/post/%s/revision/%d", profile.Service, profile.Id, postId, revisionId)
	resp, err := f.GetResult(ctx, url, headers, &response)

	if err != nil {
		return types.Result[Response]{Err: err}
	} else if resp.IsError() {
		return types.Result[Response]{Err: fmt.Errorf("error fetching user '%s' posts %s",
			profile.Name, resp.Status())}
	}

	return types.Result[Response]{Data: Response{
		Post:   response.Post,
		Images: response.Images,
		Videos: response.Videos,
	}}
}
