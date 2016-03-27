package main

import (
	"log"

	"github.com/jordanpotter/remote-backup/backup/google"
)

func main() {
	err := google.Backup("vendor", "backups")
	if err != nil {
		log.Printf("Unexpected backup error: %v", err)
	}
}
