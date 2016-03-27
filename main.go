package main

import "github.com/jordanpotter/remote-backup/backup/google"

func main() {
	google.Backup(".", "backups")
}
