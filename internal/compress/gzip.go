package compress

import (
	"compress/gzip"
	"io"
)

const compressionLevel = gzip.BestCompression

func Gzip(src io.ReadCloser) (io.ReadCloser, <-chan error) {
	errc := make(chan error, 2)
	r, w := io.Pipe()

	go func() {
		defer src.Close()
		defer w.Close()
		defer close(errc)

		gw, err := gzip.NewWriterLevel(w, compressionLevel)
		if err != nil {
			errc <- err
			return
		}
		defer gw.Close()

		_, err = io.Copy(gw, src)
		errc <- err
	}()

	return r, errc
}
