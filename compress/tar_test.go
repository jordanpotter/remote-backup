package compress

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/stretchr/testify/require"
)

const (
	testTarDirPrefix = "test_tar"
)

type testFile struct {
	path    string
	content string
}

func TestTar(t *testing.T) {
	testFiles := []testFile{
		{"a.txt", "a hello"},
		{"b/c.txt", "b/c hello"},
		{"b/d.txt", "b/d hello"},
		{"b/e/f.txt", "b/e/f hello"},
	}

	originalPath, err := ioutil.TempDir("", testTarDirPrefix)
	require.NoError(t, err)

	err = createTestFiles(originalPath, testFiles)
	require.NoError(t, err)

	var tarred bytes.Buffer
	err = Tar(originalPath, utils.NopWriteCloser(&tarred))
	require.NoError(t, err)

	untarPath, err := ioutil.TempDir("", testTarDirPrefix)
	require.NoError(t, err)

	err = Untar(untarPath, ioutil.NopCloser(&tarred))
	require.NoError(t, err)

	err = verifyTestFilesMatch(originalPath, untarPath, testFiles)
	require.NoError(t, err)
}

func createTestFiles(rootPath string, files []testFile) error {
	for _, f := range files {
		dirPath := filepath.Dir(f.path)
		if dirPath != f.path {
			absDirPath := filepath.Join(rootPath, dirPath)
			if err := os.MkdirAll(absDirPath, 0777); err != nil {
				return err
			}
		}

		absPath := filepath.Join(rootPath, f.path)
		if err := ioutil.WriteFile(absPath, []byte(f.content), 0666); err != nil {
			return err
		}
	}
	return nil
}

func verifyTestFilesMatch(originalPath, untarPath string, files []testFile) error {
	for _, f := range files {
		absOriginalPath := filepath.Join(originalPath, f.path)

		originalInfo, err := os.Stat(absOriginalPath)
		if err != nil {
			return err
		}

		originalBytes, err := ioutil.ReadFile(absOriginalPath)
		if err != nil {
			return err
		}

		absUntarPath := filepath.Join(untarPath, f.path)

		untarInfo, err := os.Stat(absUntarPath)
		if err != nil {
			return err
		}

		untarBytes, err := ioutil.ReadFile(absUntarPath)
		if err != nil {
			return err
		}

		if originalInfo.Mode() != untarInfo.Mode() {
			return fmt.Errorf("mismatched file permissions for file %q, expected %o got %o", f.path, originalInfo.Mode(), untarInfo.Mode())
		}

		if !bytes.Equal(originalBytes, untarBytes) {
			return fmt.Errorf("mismatched content for file %q", f.path)
		}
	}
	return nil
}
