package google

import (
	"bytes"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"

	"github.com/jordanpotter/remote-backup/internal/tar"
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

	buf := new(bytes.Buffer)
	err = tar.TarDir(path, buf)
	if err != nil {
		log.Fatalf("Unable to tar directory: %v", err)
	}

	log.Println("TODO: perform backup to Google Cloud Storage")
	return nil
}
