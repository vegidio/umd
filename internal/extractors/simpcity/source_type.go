package simpcity

// SourceThread represents a thread source type.
type SourceThread struct {
	id string
}

func (s SourceThread) Type() string {
	return "Thread"
}

func (s SourceThread) Name() string {
	return s.id
}
