package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func CTR(secret string, r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	hashedSecret := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(hashedSecret[:])
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
