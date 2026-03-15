package coomer

// serviceSource provides access to the service name.
type serviceSource interface {
	ServiceName() string
}

// SourceUser represents a user source type.
type SourceUser struct {
	Service string

	name string
}

func (s SourceUser) Type() string        { return "User" }
func (s SourceUser) Name() string        { return s.name }
func (s SourceUser) ServiceName() string { return s.Service }

// SourcePost represents a post source type.
type SourcePost struct {
	Service string
	Id      string

	name string
}

func (s SourcePost) Type() string        { return "Post" }
func (s SourcePost) Name() string        { return s.name }
func (s SourcePost) ServiceName() string { return s.Service }
