package compress

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jordanpotter/remote-backup/utils"
	"github.com/pkg/errors"
)

// Tar will write a tarball to w of the directory specified by rootPath.
func Tar(rootPath string, w io.WriteCloser) error {
	defer utils.MustClose(w)

	tw := tar.NewWriter(w)
	defer utils.MustClose(tw)

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

		isSymlink := info.Mode()&os.ModeSymlink != 0
		if isSymlink {
			log.Printf("Ignoring symlink %q", path)
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return errors.Wrapf(err, "failed to create info header for file %q", path)
		}

		header.Name, err = filepath.Rel(rootPath, path)
		if err != nil {
			return errors.Wrapf(err, "failed to determine relative path for file %q", path)
		}

		err = tw.WriteHeader(header)
		if err != nil {
			return errors.Wrapf(err, "failed to write header for file %q", path)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %q", path)
		}
		defer utils.MustClose(file)

		_, err = io.Copy(tw, file)
		return errors.Wrapf(err, "failed to send file %q to tar writer", path)
	}
}

// Untar will extract the given tarball data in r to the directory specified by rootPath.
func Untar(rootPath string, r io.ReadCloser) error {
	defer utils.MustClose(r)

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
		err = untarFileHandler(path, info, tr)
		if err != nil {
			return errors.Wrapf(err, "failed to handle file %q", path)
		}
	}
	return nil
}

func untarFileHandler(path string, info os.FileInfo, tr io.Reader) error {
	if info.IsDir() {
		err := os.MkdirAll(path, info.Mode())
		return errors.Wrapf(err, "failed to create directory %q", path)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return errors.Wrapf(err, "failed to open file %q", path)
	}
	defer utils.MustClose(file)

	_, err = io.Copy(file, tr)
	return errors.Wrapf(err, "failed to write to file %q", path)
}
