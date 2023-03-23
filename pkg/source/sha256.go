package source

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

func downloadFile(client Downloader, uri string) (string, error) {
	return client.DownloadFileWithName(uri, fmt.Sprintf("%d", time.Now().Unix()))
}

func getSha256ForAsset(client Downloader, uri string) (string, error) {
	file, err := downloadFile(client, uri)
	if err != nil {
		return "", err
	}

	defer os.Remove(file)
	sha256, err := getSha256(file)
	if err != nil {
		return "", err
	}

	return sha256, nil
}

func getSha256(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
