package cyberdrop

// SourceMedia represents a media source type.
type SourceMedia struct {
	id string
}

func (s SourceMedia) Type() string {
	return "Media"
}

func (s SourceMedia) Name() string {
	return s.id
}

// SourceAlbum represents an image source type.
type SourceAlbum struct {
	id string
}

func (s SourceAlbum) Type() string {
	return "Album"
}

func (s SourceAlbum) Name() string {
	return s.id
}
