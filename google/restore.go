package google

import (
	"context"
	"io"

	"github.com/jordanpotter/remote-backup/compress"
	"github.com/jordanpotter/remote-backup/crypto"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Restore will decrypt and decompress the most recent backup from the bucket
// specified and save it to the given path.
func Restore(projectID, bucket, path, secret string) error {
	ctrReader, googleWriter := io.Pipe()
	gzipReader, ctrWriter := io.Pipe()
	tarReader, gzipWriter := io.Pipe()

	group, context := errgroup.WithContext(context.Background())

	group.Go(func() error {
		err := downloadFromBucket(context, projectID, bucket, googleWriter)
		return errors.Wrapf(err, "failed to download file from bucket %q", bucket)
	})

	group.Go(func() error {
		err := crypto.CTRDecrypt(secret, ctrReader, ctrWriter)
		return errors.Wrap(err, "failed to decrypt")
	})

	group.Go(func() error {
		err := compress.Gunzip(gzipReader, gzipWriter)
		return errors.Wrap(err, "failed to decompress")
	})

	group.Go(func() error {
		err := compress.Untar(path, tarReader)
		return errors.Wrapf(err, "failed to untar to %q", path)
	})

	err := group.Wait()
	return errors.Wrap(err, "unable to restore backup")
}
