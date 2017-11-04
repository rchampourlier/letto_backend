package adapters

import (
	"os"
	"strings"

	"github.com/spf13/afero"
)

// AferoFs provides an adapter that stores workflows through an
// Afero Fs wrapper.
type AferoFs struct {
	fs      afero.Fs
	rootDir string
}

// NewAferoFs creates a new adapter.
func NewAferoFs(fs afero.Fs, rootDir string) AferoFs {
	return AferoFs{fs, rootDir}
}

// ListObjectPaths returns the paths to the objects stored
// in the underlying storage.
func (a *AferoFs) ListObjectPaths() ([]string, error) {
	paths := make([]string, 0)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		// If `err` is not nil, there had been an error walking the directory.
		// We return the error to stop the processing.
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath := strings.Replace(path, a.rootDir+"/", "", 1)
			// DEBT: it's dangerous to use string replacement to remove
			//       the parent `rootDir` from the path
			paths = append(paths, relativePath)
		}
		return nil
	}
	err := afero.Walk(a.fs, a.rootDir, walkFunc)
	if err != nil {
		return paths, err
	}
	return paths, nil
}

// CreateObject mocks the creation of an object with the
// specified properties.
func (a *AferoFs) CreateObject(path string, content string) (err error) {
	panic("Not implemented")
}

// Count returns the number of objects stored.
func (a *AferoFs) Count() (count int, err error) {
	panic("Not implemented")
}

// Reset deletes all stored objects and resets the adapter
// in its initial state.
func (a *AferoFs) Reset() {
	panic("Not implemented")
}
