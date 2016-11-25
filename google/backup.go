package google

import (
	"context"
	"io"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/jordanpotter/remote-backup/internal/compress"
	"github.com/jordanpotter/remote-backup/internal/crypto"
	"github.com/pkg/errors"
)

func Backup(projectID, bucket, path, secret string) error {
	gzipReader, tarWriter := io.Pipe()
	ctrReader, gzipWriter := io.Pipe()
	googleReader, ctrWriter := io.Pipe()

	group, context := errgroup.WithContext(context.Background())

	group.Go(func() error {
		err := compress.Tar(path, tarWriter)
		return errors.Wrapf(err, "failed to tar %q", path)
	})

	group.Go(func() error {
		err := compress.Gzip(gzipReader, gzipWriter)
		return errors.Wrap(err, "failed to gzip")
	})

	group.Go(func() error {
		err := crypto.CTREncrypt(secret, ctrReader, ctrWriter)
		return errors.Wrap(err, "failed to encrypt")
	})

	group.Go(func() error {
		now := time.Now().UTC().Format(time.RFC3339)
		err := uploadToBucket(context, projectID, bucket, now, googleReader)
		return errors.Wrapf(err, "failed to updated to bucket %q", bucket)
	})

	err := group.Wait()
	return errors.Wrap(err, "unable to create backup")
}
