package crypto

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCTR(t *testing.T) {
	secret := "secret"
	payload := "test-payload"

	var original bytes.Buffer
	_, err := original.WriteString(payload)
	require.NoError(t, err)

	var encrypted bytes.Buffer
	err = CTREncrypt(secret, ioutil.NopCloser(&original), utils.NopWriteCloser(&encrypted))
	require.NoError(t, err)

	var decrypted bytes.Buffer
	err = CTRDecrypt(secret, ioutil.NopCloser(&encrypted), utils.NopWriteCloser(&decrypted))
	require.NoError(t, err)

	assert.Equal(t, payload, decrypted.String())
}
