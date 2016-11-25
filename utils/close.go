package utils

import (
	"io"
	"log"
)

// MustClose attempts to close the provided closer and exits the program on failure.
func MustClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Fatalln("Unexpected error while closing: %v", err)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping the provided writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}
