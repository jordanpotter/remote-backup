package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func TarDir(path string, w io.Writer) error {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory path", path)
	}

	t := tar.NewWriter(w)

	handleFile := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if path != "" {
			header.Name = filepath.Join(path, filePath)
		}

		err = t.WriteHeader(header)
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

		_, err = io.Copy(t, file)
		return err
	}

	return filepath.Walk(path, handleFile)
}
