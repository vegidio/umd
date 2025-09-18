package fapello

// SourcePost represents a post source type.
type SourcePost struct {
	Id string

	name string
}

func (s SourcePost) Type() string {
	return "Post"
}

func (s SourcePost) Name() string {
	return s.name
}

// SourceModel represents a model source type.
type SourceModel struct {
	name string
}

func (s SourceModel) Type() string {
	return "Model"
}

func (s SourceModel) Name() string {
	return s.name
}
