package chunk

// Config controls how policy text is split for retrieval.
type Config struct {
	MinTokens     int
	MaxTokens     int
	OverlapTokens int
}

// DefaultConfig returns chunk settings aligned with the indexing strategy.
func DefaultConfig() Config {
	return Config{
		MinTokens:     500,
		MaxTokens:     800,
		OverlapTokens: 75,
	}
}
