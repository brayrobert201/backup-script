#!/bin/bash

# Define the source and destination directories
SOURCE_DIR="/home/"
DEST_DIR="/mnt/d/WSL/"

# Check if zstd is installed
if ! command -v zstd &> /dev/null
then
    echo "zstd could not be found, please install it."
    exit 1
fi

# Create a timestamp
TIMESTAMP=$(date +"%Y%m%d%H%M%S")

# Create the backup filename with the timestamp
BACKUP_FILE="backup-${TIMESTAMP}.tar.zstd"

# Check if destination directory exists, create if not
if [ ! -d "$DEST_DIR" ]; then
    mkdir -p "$DEST_DIR"
fi

check_integrity() {
    echo "Checking integrity of $1"
    if zstd -t $1; then
        echo "Integrity check passed for $1"
    else
        echo "Integrity check failed for $1"
        exit 3
    fi
}

# Navigate to the source directory
cd $SOURCE_DIR

# Use tar combined with zstd to create a compressed backup
if tar cf - . | pzstd -9 > "${DEST_DIR}${BACKUP_FILE}"; then
    echo "Backup created successfully: ${DEST_DIR}${BACKUP_FILE}"
    check_integrity "${DEST_DIR}${BACKUP_FILE}"
else
    echo "Backup failed."
    exit 2
fi

# Optional: Cleanup older backups
# Uncomment the line below and adjust DAYS to your preference
 find $DEST_DIR -name 'backup-*.tar.zstd' -mtime +7 -delete

