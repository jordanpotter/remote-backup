package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/pkg/errors"
)

func CTREncrypt(secret string, r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	iv, err := generateIV()
	if err != nil {
		return errors.Wrap(err, "failed to generate IV")
	}

	w.Write(iv[:])
	return ctr(secret, iv, r, w)
}

func CTRDecrypt(secret string, r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()

	iv, err := readIV(r)
	if err != nil {
		return errors.Wrap(err, "failed to read IV")
	}

	return ctr(secret, iv, r, w)
}

func ctr(secret string, iv []byte, r io.ReadCloser, w io.WriteCloser) error {
	block, err := aes.NewCipher(hashedSecret(secret))
	if err != nil {
		return errors.Wrap(err, "failed to create AES cipher")
	}

	cipher.NewCTR(block, iv[:])
	stream := cipher.NewOFB(block, iv[:])
	sr := &cipher.StreamReader{S: stream, R: r}
	_, err = io.Copy(w, sr)
	return errors.Wrap(err, "failed to perform cipher on data")
}
