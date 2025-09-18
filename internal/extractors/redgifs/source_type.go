package redgifs

// SourceVideo represents a video source type.
type SourceVideo struct {
	name string
}

func (s SourceVideo) Type() string {
	return "Video"
}

func (s SourceVideo) Name() string {
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
