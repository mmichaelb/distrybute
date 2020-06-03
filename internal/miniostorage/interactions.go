package miniostorage

import (
	"github.com/minio/minio-go"
	"github.com/mmichaelb/distrybute/internal"
	"io"
)

// PutFile is the minio client based implementation of the Storage interface PutFile function.
func (instance *Instance) PutFile(id string, reader io.Reader) (err error) {
	_, err = instance.client.PutObject(instance.bucketName, id, reader, -1, minio.PutObjectOptions{})
	return
}

// GetFile is the minio client based implementation of the Storage interface GetFile function.
func (instance *Instance) GetFile(id string) (distrybute.ReadCloseSeeker, error) {
	obj, err := instance.client.GetObject(instance.bucketName, id, minio.GetObjectOptions{})
	return obj, err
}

// DeleteFile is the minio client based implementation of the Storage interface DeleteFile function.
func (instance *Instance) DeleteFile(id string) (err error) {
	err = instance.client.RemoveObject(instance.bucketName, id)
	return
}

// DeleteMultipleFiles is the minio client based implementation of the Storage interface DeleteMultipleFiles function.
func (instance *Instance) DeleteMultipleFiles(ids ...string) (errors []error) {
	objectsChan := make(chan string, len(ids))
	for _, id := range ids {
		objectsChan <- id
	}
	defer close(objectsChan)
	errorChan := instance.client.RemoveObjects(instance.bucketName, objectsChan)
	// Print errors received from RemoveObjects API
	for err := range errorChan {
		if errors == nil {
			errors = make([]error, 1)
			errors[0] = err.Err
		} else {
			errors = append(errors, err.Err)
		}
	}
	return
}
