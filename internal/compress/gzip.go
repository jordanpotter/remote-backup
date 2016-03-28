package compress

import (
	"compress/gzip"
	"io"
)

const compressionLevel = gzip.BestCompression

func Gzip(r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	gw, err := gzip.NewWriterLevel(w, compressionLevel)
	if err != nil {
		return err
	}
	defer gw.Close()

	_, err = io.Copy(gw, r)
	return err
}
