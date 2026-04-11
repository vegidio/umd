package shared

type Download struct {
	Url       string
	FilePath  string
	Error     error
	IsSuccess bool
	Hash      string
}
