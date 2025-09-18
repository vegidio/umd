package types

type ExtractorType int

const (
	// Generic represents a generic extractor type.
	Generic ExtractorType = iota
	// Coomer represents the Coomer (coomer.st) extractor type.
	Coomer
	// Fapello represents the Fapello (fapello.com) extractor type.
	Fapello
	// Imaglr represents the Imaglr (imaglr.com) extractor type.
	Imaglr
	// Reddit represents the Reddit (reddit.com) extractor type.
	Reddit
	// RedGifs the RedGifs (redgifs.com) extractor type.
	RedGifs
	// Kemono the Kemono (kemono.cr) extractor type.
	Kemono
)

func (e ExtractorType) String() string {
	switch e {
	case Generic:
		return "Generic"
	case Coomer:
		return "Coomer"
	case Fapello:
		return "Fapello"
	case Imaglr:
		return "Imaglr"
	case Kemono:
		return "Kemono"
	case Reddit:
		return "Reddit"
	case RedGifs:
		return "RedGifs"
	}

	return "Unknown"
}
