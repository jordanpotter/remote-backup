package google

import (
	"io"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"

	"github.com/jordanpotter/remote-backup/internal/compress"
	"github.com/jordanpotter/remote-backup/internal/encrypt"
)

func Backup(path, bucket string) error {
	client, err := google.DefaultClient(oauth2.NoContext, storage.DevstorageFullControlScope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}

	_, err = storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	f, err := ioutil.TempFile("", "test")
	if err != nil {
		return err
	}

	err = processFiles(path, f)
	if err != nil {
		log.Fatalf("Unable to process files: %v", err)
	}
	return nil
}

func processFiles(path string, w io.WriteCloser) error {
	var wg sync.WaitGroup
	errc := make(chan error)

	gr, tw := io.Pipe()
	cr, gw := io.Pipe()

	wg.Add(1)
	go func() {
		errc <- compress.Tar(path, tw)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errc <- compress.Gzip(gr, gw)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errc <- encrypt.CTR("example key 1234", cr, w)
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
