package miniostorage

import (
	"errors"
	"github.com/minio/minio-go"
	"io"
)

// Instance contains the functions needed to implement the Storage interface declares inside the
// root package. For instantiation, use the New function.
type Instance struct {
	client     *minio.Client
	bucketName string
}

// New instantiates a new instance and checks if the given bucket is already existent and if not,
// creates it. Additionally, it checks if all parameters are valid.
func New(client *minio.Client, bucketName string, bucketLocation string) (*Instance, error) {
	if client == nil {
		return nil, errors.New("given minio client instance cannot be nil")
	}
	if bucketName == "" {
		return nil, errors.New("given bucket name cannot be empty")
	}
	if ok, err := client.BucketExists(bucketName); err != nil {
		return nil, err
	} else if !ok {
		if err := client.MakeBucket(bucketName, bucketLocation); err != nil {
			return nil, err
		}
	}
	return &Instance{client: client, bucketName: bucketName}, nil
}

// PutFile is the minio client based implentation of the Storage interface PutFile function.
func (instance *Instance) PutFile(id string, reader io.Reader) (err error) {
	_, err = instance.client.PutObject(instance.bucketName, id, reader, -1, minio.PutObjectOptions{})
	return
}
