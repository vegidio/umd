package erome

// SourceAlbum represents an album source type.
type SourceAlbum struct {
	Id string
}

func (s SourceAlbum) Type() string {
	return "Album"
}

func (s SourceAlbum) Name() string {
	return s.Id
}
