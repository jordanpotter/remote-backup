package main

import (
	"encoding/hex"
	"log"
	"flag"

	"github.com/jordanpotter/remote-backup/backup/google"
)

var (
	projectID string
	bucketName string
	path string
	secretKeyHex string
)

func init() {
	flag.StringVar(&projectID, "project", "", "Google Cloud project id")
	flag.StringVar(&bucketName, "bucket", "", "bucket name to backup to")
	flag.StringVar(&path, "path", "", "directory path to backup")
	flag.StringVar(&secretKeyHex, "key", "", "secret encryption key in hexadecimal")
	flag.Parse()
}

func main() {
	verifyFlags()

	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		log.Fatalf("Unexpected backup error: %v", err)
	}

	err = google.Backup(projectID, bucketName, path, secretKey)
	if err != nil {
		log.Fatalf("Unexpected backup error: %v", err)
	}
}

func verifyFlags() {
	if projectID == "" {
		log.Fatalln("Must specify project id")
	} else if bucketName == "" {
		log.Fatalln("Must specify bucket name")
	} else if path == "" {
		log.Fatalln("Must specify path")
	} else if secretKeyHex == "" {
		log.Fatalln("Must specify secret key")
	}
}
