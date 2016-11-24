package google

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func uploadToBucket(projectID, bucket, filename string, r io.ReadCloser) error {
	defer r.Close()

	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	s, err := storage.New(c)
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}

	err = ensureBucketExists(s, projectID, bucket)
	if err != nil {
		return fmt.Errorf("failed to ensure bucket %q exists: %v", bucket, err)
	}

	o := &storage.Object{Name: filename}
	_, err = s.Objects.Insert(bucket, o).Media(r).Do()
	if err != nil {
		return fmt.Errorf("failed to upload object to bucket %q: %v", bucket, err)
	}
	return nil
}

func ensureBucketExists(service *storage.Service, projectID, bucket string) error {
	_, err := service.Buckets.Get(bucket).Do()
	if err == nil {
		return nil
	}

	b := &storage.Bucket{Name: bucket, StorageClass: "NEARLINE"}
	_, err = service.Buckets.Insert(projectID, b).Do()
	if err != nil {
		return fmt.Errorf("failed to create bucket %q: %v", bucket, err)
	}
	return nil
}
