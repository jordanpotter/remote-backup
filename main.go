package main

import (
	"flag"
	"log"

	"github.com/jordanpotter/remote-backup/google"
)

var (
	backup    bool
	restore   bool
	projectID string
	bucket    string
	path      string
	secret    string
)

func init() {
	flag.BoolVar(&backup, "backup", false, "backup data")
	flag.BoolVar(&restore, "restore", false, "restore data")
	flag.StringVar(&projectID, "project", "", "Google Cloud project id")
	flag.StringVar(&bucket, "bucket", "", "bucket to store backups")
	flag.StringVar(&path, "path", "", "path to directory")
	flag.StringVar(&secret, "secret", "", "secret for encryption/decryption")
	flag.Parse()
}

func main() {
	verifyFlags()

	if backup {
		if err := google.Backup(projectID, bucket, path, secret); err != nil {
			log.Fatalf("Unexpected backup error: %v", err)
		}
	}

	if restore {
		if err := google.Restore(projectID, bucket, path, secret); err != nil {
			log.Fatalf("Unexpected restore error: %v", err)
		}
	}
}

func verifyFlags() {
	if !backup && !restore {
		log.Fatalln("Must specify backup or restore")
	} else if backup && restore {
		log.Fatalln("Cannot specify both backup and restore")
	} else if projectID == "" {
		log.Fatalln("Must specify project")
	} else if bucket == "" {
		log.Fatalln("Must specify bucket")
	} else if path == "" {
		log.Fatalln("Must specify path")
	} else if secret == "" {
		log.Fatalln("Must specify encryption/decryption secret")
	}
}
