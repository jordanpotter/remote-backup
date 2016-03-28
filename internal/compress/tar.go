package compress

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Tar(path string, w io.WriteCloser) error {
	defer w.Close()

	tw := tar.NewWriter(w)
	defer tw.Close()

	info, err := os.Stat(path)
	if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory path", path)
	} else {
		return filepath.Walk(path, getFileTarHandler(path, tw))
	}
}

func getFileTarHandler(rootPath string, tw *tar.Writer) filepath.WalkFunc {
	return func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name = filepath.Join(rootPath, filePath)
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	}
}
