package types

type MediaType int

const (
	Image MediaType = iota
	Video
	Unknown
)

func (m MediaType) String() string {
	switch m {
	case Image:
		return "Image"
	case Video:
		return "Video"
	default:
		return "Unknown"
	}
}
