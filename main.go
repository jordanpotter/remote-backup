package main

import (
	"encoding/hex"
	"log"

	"github.com/jordanpotter/remote-backup/backup/google"
)

func main() {
	secretKey, err := hex.DecodeString("ABCDABCDABCDABCDABCDABCDABCDABCD")
	if err != nil {
		log.Fatalf("Unexpected backup error: %v", err)
	}

	err = google.Backup("vendor", secretKey, "backups")
	if err != nil {
		log.Fatalf("Unexpected backup error: %v", err)
	}
}
