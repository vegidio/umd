package imaglr

// SourcePost represents a post source type.
type SourcePost struct {
	name string
}

func (s SourcePost) Type() string {
	return "Post"
}

func (s SourcePost) Name() string {
	return s.name
}
