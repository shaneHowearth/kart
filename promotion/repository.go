package promotion

type Store interface {
	GetCodeFileMatchCount(string) (int, error)
	AddCodeFileMatchCount(string, int) error
	InitialiseDataStore() error
}
