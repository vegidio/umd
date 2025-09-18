package reddit

// SourceSubmission represents a submission source type.
type SourceSubmission struct {
	Id string

	name string
}

func (s SourceSubmission) Type() string {
	return "Submission"
}

func (s SourceSubmission) Name() string {
	return s.name
}

// SourceUser represents a user source type.
type SourceUser struct {
	name string
}

func (s SourceUser) Type() string {
	return "User"
}

func (s SourceUser) Name() string {
	return s.name
}

// SourceSubreddit represents a subreddit source type.
type SourceSubreddit struct {
	name string
}

func (s SourceSubreddit) Type() string {
	return "Subreddit"
}

func (s SourceSubreddit) Name() string {
	return s.name
}
