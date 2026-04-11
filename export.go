package umd

import "github.com/vegidio/umd/internal/types"

type Response = types.Response
type ExtractorType = types.ExtractorType
type Media = types.Media
type MediaType = types.MediaType
type Metadata = types.Metadata

// GetMediaType returns the MediaType for a given file extension (without leading dot).
var GetMediaType = types.GetType

const (
	Generic   = types.Generic
	Bunkr     = types.Bunkr
	Coomer    = types.Coomer
	Cyberdrop = types.Cyberdrop
	Erome     = types.Erome
	Fapello   = types.Fapello
	Imaglr    = types.Imaglr
	JpgFish   = types.JpgFish
	Kemono    = types.Kemono
	Reddit    = types.Reddit
	RedGifs   = types.RedGifs
	Saint     = types.Saint
	SimpCity  = types.SimpCity
)
