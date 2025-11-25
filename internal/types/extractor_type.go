package types

type ExtractorType int

const (
	// Generic represents a generic extractor type.
	Generic ExtractorType = iota
	// Bunkr represents the Bunkr (bunkr.cr) extractor type.
	Bunkr
	// Coomer represents the Coomer (coomer.st) extractor type.
	Coomer
	// Cyberdrop represents the Cyberdrop (cyberdrop.cr) extractor type.
	Cyberdrop
	// Erome represents the Erome (erome.com) extractor type.
	Erome
	// Fapello represents the Fapello (fapello.com) extractor type.
	Fapello
	// Imaglr represents the Imaglr (imaglr.com) extractor type.
	Imaglr
	// JpgFish represents the JpgFish (jpg6.su) extractor type.
	JpgFish
	// Kemono represents the Kemono (kemono.cr) extractor type.
	Kemono
	// Reddit represents the Reddit (reddit.com) extractor type.
	Reddit
	// RedGifs represents the RedGifs (redgifs.com) extractor type.
	RedGifs
	// Saint represents the Saint (saint2.su) extractor type.
	Saint
	// SimpCity represents the SimpCity (simpcity.cr) extractor type.
	SimpCity
)

func (e ExtractorType) String() string {
	switch e {
	case Generic:
		return "Generic"
	case Bunkr:
		return "Bunkr"
	case Coomer:
		return "Coomer"
	case Cyberdrop:
		return "Cyberdrop"
	case Erome:
		return "Erome"
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
	case SimpCity:
		return "SimpCity"
	}

	return "Unknown"
}
