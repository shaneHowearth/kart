package promotion

type Store interface {
	GetCodeFileMatchCounts([]string) (map[string]CacheResult, error)
	AddCodeFileMatchCounts(map[string]int) error
	InitialiseDataStore() error
}

type CacheResult struct {
	MatchCount int
	Found      bool
}
