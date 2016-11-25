package google

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func uploadToBucket(context context.Context, projectID, bucket, filename string, r io.ReadCloser) error {
	defer r.Close()

	c, err := google.DefaultClient(context, storage.DevstorageFullControlScope)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	s, err := storage.New(c)
	if err != nil {
		return errors.Wrap(err, "failed to create storage service")
	}

	err = ensureBucketExists(s, projectID, bucket)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure bucket %q exists", bucket)
	}

	o := &storage.Object{Name: filename}
	_, err = s.Objects.Insert(bucket, o).Media(r).Do()
	return errors.Wrapf(err, "failed to upload to bucket %q", bucket)
}

func ensureBucketExists(service *storage.Service, projectID, bucket string) error {
	_, err := service.Buckets.Get(bucket).Do()
	if err == nil {
		return nil
	}

	b := &storage.Bucket{Name: bucket, StorageClass: "NEARLINE"}
	_, err = service.Buckets.Insert(projectID, b).Do()
	return errors.Wrapf(err, "failed to create bucket %q", bucket)
}
