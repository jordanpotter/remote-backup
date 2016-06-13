package google

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func uploadToBucket(projectID, bucketName, filename string, r io.ReadCloser) error {
	defer r.Close()

	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	s, err := storage.New(c)
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}

	err = ensureBucketExists(s, projectID, bucketName)
	if err != nil {
		return fmt.Errorf("failed to ensure bucket %q exists: %v", bucketName, err)
	}

	o := &storage.Object{Name: filename}
	_, err = s.Objects.Insert(bucketName, o).Media(r).Do()
	if err != nil {
		return fmt.Errorf("failed to upload object to bucket %q: %v", bucketName, err)
	}
	return nil
}

func ensureBucketExists(service *storage.Service, projectID, bucketName string) error {
	_, err := service.Buckets.Get(bucketName).Do()
	if err == nil {
		return nil
	}

	b := &storage.Bucket{Name: bucketName, StorageClass: "NEARLINE"}
	_, err = service.Buckets.Insert(projectID, b).Do()
	if err != nil {
		return fmt.Errorf("failed to create bucket %q: %v", bucketName, err)
	}
	return nil
}
