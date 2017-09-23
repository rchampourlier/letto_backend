package adapters

import (
	"fmt"
)

// S3 provides an interface with S3 to be used as a storage
// facility
type S3 struct {
	bucket string
}

// NewS3 creates a new S3Adapter
func NewS3(bucket string) S3 {
	return S3{bucket}
}

// CreateObject creates the specified object on AWS S3.
func (s3 *S3) CreateObject(path string, source string) (err error) {
	// TODO
	fmt.Println("Not implemented")
	return nil
}
