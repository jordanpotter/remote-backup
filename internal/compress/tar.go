package compress

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Tar(path string) (io.ReadCloser, <-chan error) {
	errc := make(chan error, 1)
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		defer close(errc)

		tw := tar.NewWriter(w)
		defer tw.Close()

		info, err := os.Stat(path)
		if err != nil {
			errc <- err
		} else if !info.IsDir() {
			errc <- fmt.Errorf("%s is not a directory path", path)
		} else {
			errc <- filepath.Walk(path, getFileTarHandler(path, tw))
		}
	}()

	return r, errc
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
