package google

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func Backup() error {
	client, err := google.DefaultClient(oauth2.NoContext, storage.DevstorageFullControlScope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}

	_, err = storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	log.Println("TODO: perform backup to Google Cloud Storage")
	return nil
}
