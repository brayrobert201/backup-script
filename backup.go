package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	sourceDir = "/home/"
	destDir   = "/mnt/d/WSL/"
)

func main() {
	// Check if zstd is installed
	_, err := exec.LookPath("zstd")
	if err != nil {
		fmt.Println("zstd could not be found, please install it.")
		os.Exit(1)
	}

	// Create a timestamp
	timestamp := time.Now().Format("20060102150405")

	// Create the backup filename with the timestamp
	backupFile := fmt.Sprintf("backup-%s.tar.zstd", timestamp)

	// Check if destination directory exists, create if not
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		os.MkdirAll(destDir, os.ModePerm)
	}

	// Navigate to the source directory
	os.Chdir(sourceDir)

	// Use tar combined with zstd to create a compressed backup
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tar cf - . | zstd -9 > %s", filepath.Join(destDir, backupFile)))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Backup failed.")
		os.Exit(2)
	} else {
		fmt.Printf("Backup created successfully: %s\n", filepath.Join(destDir, backupFile))
		checkIntegrity(filepath.Join(destDir, backupFile))
	}

	// Optional: Cleanup older backups
	// Uncomment the line below and adjust DAYS to your preference
	// find $DEST_DIR -name 'backup-*.tar.zstd' -mtime +7 -delete
}

func checkIntegrity(file string) {
	fmt.Printf("Checking integrity of %s\n", file)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("zstd -t %s", file))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Integrity check failed for %s\n", file)
		os.Exit(3)
	} else {
		fmt.Printf("Integrity check passed for %s\n", file)
	}
}
