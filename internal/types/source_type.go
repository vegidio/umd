package types

// SourceType represents a source with methods to retrieve its type and name.
type SourceType interface {
	// Type returns the type of the source.
	Type() string
	// Name returns the name of the source.
	Name() string
}
