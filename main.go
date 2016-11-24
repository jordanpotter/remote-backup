package main

import (
	"flag"
	"log"

	"github.com/jordanpotter/remote-backup/backup/google"
)

var (
	projectID string
	bucket    string
	path      string
	secret    string
)

func init() {
	flag.StringVar(&projectID, "project", "", "Google Cloud project id")
	flag.StringVar(&bucket, "bucket", "", "bucket to store backups")
	flag.StringVar(&path, "path", "", "path to directory")
	flag.StringVar(&secret, "secret", "", "secret for encryption/decryption")
	flag.Parse()
}

func main() {
	verifyFlags()

	err := google.Backup(projectID, bucket, path, secret)
	if err != nil {
		log.Fatalf("Unexpected backup error: %v", err)
	}
}

func verifyFlags() {
	if projectID == "" {
		log.Fatalln("Must specify project")
	} else if bucket == "" {
		log.Fatalln("Must specify bucket")
	} else if path == "" {
		log.Fatalln("Must specify path")
	} else if secret == "" {
		log.Fatalln("Must specify encryption/decryption secret")
	}
}
