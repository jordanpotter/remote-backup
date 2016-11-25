package crypto

import (
	"crypto/aes"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	return iv, errors.Wrapf(err, "unable to generate %d random bytes", aes.BlockSize)
}

func readIV(r io.Reader) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(r, iv)
	return iv, errors.Wrapf(err, "unable to read %d bytes", aes.BlockSize)
}
