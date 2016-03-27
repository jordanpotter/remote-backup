package google

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"

	"github.com/jordanpotter/remote-backup/internal/compress"
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

	r, errc := readFiles(path)
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println(len(buf))

	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func readFiles(path string) (io.ReadCloser, <-chan error) {
	var wg sync.WaitGroup
	errc := make(chan error, 2)

	wg.Add(1)
	tr, terrc := compress.Tar(path)

	wg.Add(1)
	gr, gerrc := compress.Gzip(tr)

	go func() {
		for err := range terrc {
			if err != nil {
				log.Printf("Unable to tar directory: %v", err)
				errc <- err
			}
		}
		wg.Done()
	}()

	go func() {
		for err := range gerrc {
			if err != nil {
				log.Printf("Unable to gzip directory: %v", err)
				errc <- err
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	return gr, errc
}
