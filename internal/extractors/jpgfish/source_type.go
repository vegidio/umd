package jpgfish

// SourceImage represents a image source type.
type SourceImage struct {
	id string
}

func (s SourceImage) Type() string {
	return "Image"
}

func (s SourceImage) Name() string {
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
