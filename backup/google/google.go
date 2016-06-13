package google

import (
	"io"
	"sync"
	"time"

	"github.com/jordanpotter/remote-backup/internal/compress"
	"github.com/jordanpotter/remote-backup/internal/encrypt"
)

func Backup(projectID, bucketName, path string, secretKey []byte) error {
	var wg sync.WaitGroup
	errc := make(chan error)

	gzipReader, tarWriter := io.Pipe()
	ctrReader, gzipWriter := io.Pipe()
	googleReader, ctrWriter := io.Pipe()

	wg.Add(1)
	go func() {
		errc <- compress.Tar(path, tarWriter)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errc <- compress.Gzip(gzipReader, gzipWriter)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errc <- encrypt.CTR(secretKey, ctrReader, ctrWriter)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errc <- uploadToBucket(projectID, bucketName, getFilename(), googleReader)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func getFilename() string {
	return time.Now().UTC().Format(time.RFC3339)
}
