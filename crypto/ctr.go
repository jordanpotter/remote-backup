package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/pkg/errors"
)

// CTREncrypt will encrypt the given data in r using AES in CTR mode and write it to w.
func CTREncrypt(secret string, r io.ReadCloser, w io.WriteCloser) error {
	defer utils.MustClose(r)
	defer utils.MustClose(w)

	iv, err := generateIV()
	if err != nil {
		return errors.Wrap(err, "failed to generate IV")
	}

	_, err = w.Write(iv[:])
	if err != nil {
		return errors.Wrap(err, "failed to write IV")
	}

	return ctr(secret, iv, r, w)
}

// CTRDecrypt will decrypt the given data in r using AES in CTR mode and write it to w.
func CTRDecrypt(secret string, r io.ReadCloser, w io.WriteCloser) error {
	defer utils.MustClose(r)
	defer utils.MustClose(w)

	iv, err := readIV(r)
	if err != nil {
		return errors.Wrap(err, "failed to read IV")
	}

	return ctr(secret, iv, r, w)
}

func ctr(secret string, iv []byte, r io.Reader, w io.Writer) error {
	block, err := aes.NewCipher(hashedSecret(secret))
	if err != nil {
		return errors.Wrap(err, "failed to create AES cipher")
	}

	stream := cipher.NewCTR(block, iv[:])
	sr := &cipher.StreamReader{S: stream, R: r}
	_, err = io.Copy(w, sr)
	return errors.Wrap(err, "failed to perform cipher on data")
}
