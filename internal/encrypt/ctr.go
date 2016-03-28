package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func CTR(key []byte, r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	cipher.NewCTR(block, iv[:])
	stream := cipher.NewOFB(block, iv[:])
	sr := &cipher.StreamReader{S: stream, R: r}
	_, err = io.Copy(w, sr)
	return err
}
