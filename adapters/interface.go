package adapters

// Adapter interfaces Letto with a storage provider (e.g. AWS S3).
// It is provided as an interface to enable the support of other
// providers in the future but it also enables to create a MockAdapter
// for testing.
type Adapter interface {
	CreateObject(path string, source string) (err error)
}
