package adapters

// InMemory provides an adapter that stores workflows in memory.
// This adapter can be used for development or testing purposes.
type InMemory struct {
	objects []object
}

// InMemoryWorkflow is the structure used to store a workflow
// value in memory.
type object struct {
	path    string
	content string
}

// NewInMemory creates a new Mock adapter.
func NewInMemory() InMemory {
	return InMemory{}
}

// CreateObject mocks the creation of an object with the
// specified properties.
func (a *InMemory) CreateObject(path string, content string) (err error) {
	a.objects = append(a.objects, object{path: path, content: content})
	return nil
}

// Count returns the number of objects stored.
func (a *InMemory) Count() (count int, err error) {
	return len(a.objects), nil
}

// Reset deletes all stored objects and resets the adapter
// in its initial state.
func (a *InMemory) Reset() {
	a.objects = make([]object, 0)
}
