package compress

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Tar(rootPath string, w io.WriteCloser) error {
	defer w.Close()

	tw := tar.NewWriter(w)
	defer tw.Close()

	info, err := os.Stat(rootPath)
	if err != nil {
		return errors.Wrapf(err, "failed to get stats for %q", rootPath)
	} else if !info.IsDir() {
		return errors.Errorf("%q is not a directory", rootPath)
	} else {
		err = filepath.Walk(rootPath, tarFileHandler(rootPath, tw))
		return errors.Wrapf(err, "failed to walk directory %q", rootPath)
	}
}

func tarFileHandler(rootPath string, tw *tar.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return errors.Wrap(err, "failed to create file info header")
		}

		header.Name, err = filepath.Rel(rootPath, path)
		if err != nil {
			return errors.Wrap(err, "failed to get relative path of file")
		}

		err = tw.WriteHeader(header)
		if err != nil {
			return errors.Wrap(err, "failed to write file header")
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %q", path)
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return errors.Wrap(err, "failed to send file data to tar writer")
	}
}

func Untar(rootPath string, r io.ReadCloser) error {
	defer r.Close()

	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "failed to retrieve next header")
		}

		path := filepath.Join(rootPath, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return errors.Wrapf(err, "failed to create directory %q", path)
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return errors.Wrapf(err, "failed to open file %q", path)
		}
		defer file.Close()

		if _, err = io.Copy(file, tr); err != nil {
			return errors.Wrapf(err, "failed to write to file %q", path)
		}
	}
	return nil
}
