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
	// JpgFish represents the JpgFish (jpg6.su) extractor type.
	JpgFish
	// Kemono the Kemono (kemono.cr) extractor type.
	Kemono
	// Reddit represents the Reddit (reddit.com) extractor type.
	Reddit
	// RedGifs represents the RedGifs (redgifs.com) extractor type.
	RedGifs
	// Saint represents the Saint (saint2.su) extractor type.
	Saint
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
	case JpgFish:
		return "JpgFish"
	case Kemono:
		return "Kemono"
	case Reddit:
		return "Reddit"
	case RedGifs:
		return "RedGifs"
	case Saint:
		return "Saint"
	}

	return "Unknown"
}
