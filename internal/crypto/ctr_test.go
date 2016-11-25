package crypto

import (
	"testing"

	"github.com/jordanpotter/remote-backup/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCTR(t *testing.T) {
	secret := "secret"
	payload := "test-payload"

	var original testutils.Buffer
	original.WriteString(payload)

	var encrypted testutils.Buffer
	err := CTREncrypt(secret, &original, &encrypted)
	require.NoError(t, err)

	var decrypted testutils.Buffer
	err = CTRDecrypt(secret, &encrypted, &decrypted)
	require.NoError(t, err)

	assert.Equal(t, payload, decrypted.String())
}
