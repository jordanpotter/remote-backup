package compress

import (
	"testing"

	"github.com/jordanpotter/remote-backup/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzip(t *testing.T) {
	payload := "test-payload"

	var original testutils.Buffer
	original.WriteString(payload)

	var compressed testutils.Buffer
	err := Gzip(&original, &compressed)
	require.NoError(t, err)

	var decompressed testutils.Buffer
	err = Gunzip(&compressed, &decompressed)
	require.NoError(t, err)

	assert.Equal(t, payload, decompressed.String())
}
