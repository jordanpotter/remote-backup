package google

import (
	"context"
	"io"
	"time"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/pkg/errors"

	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func uploadToBucket(context context.Context, projectID, bucket, filename string, r io.ReadCloser) error {
	defer utils.MustClose(r)

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

func downloadFromBucket(context context.Context, projectID, bucket string, w io.WriteCloser) error {
	defer utils.MustClose(w)

	c, err := google.DefaultClient(context, storage.DevstorageFullControlScope)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	s, err := storage.New(c)
	if err != nil {
		return errors.Wrap(err, "failed to create storage service")
	}

	filename, err := mostRecentFilename(s, projectID, bucket)
	if err != nil {
		return errors.Wrapf(err, "failed to find most recent file in bucket %q", bucket)
	}

	resp, err := s.Objects.Get(bucket, filename).Download()
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve file %q from bucket %q", filename, bucket)
	}
	defer utils.MustClose(resp.Body)

	_, err = io.Copy(w, resp.Body)
	return errors.Wrap(err, "failed to read response data")
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

func mostRecentFilename(service *storage.Service, projectID, bucket string) (string, error) {
	var mostRecent string
	var mostRecentCreatedTime time.Time

	objects, err := allObjectsInBucket(service, projectID, bucket)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve objects in bucket %q", bucket)
	} else if len(objects) == 0 {
		return "", errors.Errorf("no objects in bucket %q", bucket)
	}

	for _, object := range objects {
		createdTime, err := time.Parse(time.RFC3339, object.TimeCreated)
		if err != nil {
			return "", errors.Wrapf(err, "failed to parse time %q", object.TimeCreated)
		}

		if createdTime.After(mostRecentCreatedTime) {
			mostRecent = object.Name
			mostRecentCreatedTime = createdTime
		}
	}

	return mostRecent, nil
}

func allObjectsInBucket(service *storage.Service, projectID, bucket string) ([]*storage.Object, error) {
	var objects []*storage.Object

	var pageToken string
	for {
		call := service.Objects.List(bucket)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		resp, err := call.Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve files")
		}

		objects = append(objects, resp.Items...)

		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return objects, nil
}
