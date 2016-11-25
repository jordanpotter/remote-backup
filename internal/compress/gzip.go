package compress

import (
	"compress/gzip"
	"io"

	"github.com/pkg/errors"
)

const compressionLevel = gzip.DefaultCompression

func Gzip(r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	gw, err := gzip.NewWriterLevel(w, compressionLevel)
	if err != nil {
		return errors.Wrap(err, "failed to create gzip writer")
	}
	defer gw.Close()

	_, err = io.Copy(gw, r)
	return errors.Wrap(err, "failed to write gzipped data")
}

func Gunzip(r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "failed to create gzip reader")
	}
	defer gr.Close()

	_, err = io.Copy(w, gr)
	return errors.Wrap(err, "failed to read gzipped data")
}
