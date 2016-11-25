package compress

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzip(t *testing.T) {
	payload := "test-payload"

	var original bytes.Buffer
	_, err := original.WriteString(payload)
	require.NoError(t, err)

	var compressed bytes.Buffer
	err = Gzip(ioutil.NopCloser(&original), utils.NopWriteCloser(&compressed))
	require.NoError(t, err)

	var decompressed bytes.Buffer
	err = Gunzip(ioutil.NopCloser(&compressed), utils.NopWriteCloser(&decompressed))
	require.NoError(t, err)

	assert.Equal(t, payload, decompressed.String())
}
