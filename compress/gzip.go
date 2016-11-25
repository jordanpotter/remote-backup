package compress

import (
	"compress/gzip"
	"io"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/pkg/errors"
)

const compressionLevel = gzip.DefaultCompression

// Gzip will compress the given data in r using the gzip compression algorithm and write it to w.
func Gzip(r io.ReadCloser, w io.WriteCloser) error {
	defer utils.MustClose(r)
	defer utils.MustClose(w)

	gw, err := gzip.NewWriterLevel(w, compressionLevel)
	if err != nil {
		return errors.Wrap(err, "failed to create gzip writer")
	}
	defer utils.MustClose(gw)

	_, err = io.Copy(gw, r)
	return errors.Wrap(err, "failed to write gzipped data")
}

// Gunzip will decompress the given data in r using the gzip compression algorithm and write it to w.
func Gunzip(r io.ReadCloser, w io.WriteCloser) error {
	defer utils.MustClose(r)
	defer utils.MustClose(w)

	gr, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "failed to create gzip reader")
	}
	defer utils.MustClose(gr)

	_, err = io.Copy(w, gr)
	return errors.Wrap(err, "failed to read gzipped data")
}
