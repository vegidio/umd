package saint

// SourceVideo represents a video source type.
type SourceVideo struct {
	id string
}

func (s SourceVideo) Type() string {
	return "Video"
}

func (s SourceVideo) Name() string {
	return s.id
}
