package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const chunkSize = 1024 * 1024 

func uploadFileInChunks(filePath string) error {
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(filePath))
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, chunkSize)
	chunkNum := 1

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		url := fmt.Sprintf("http://localhost:8080/uploads?path=%s&chunk=%d",
			filepath.ToSlash(fileName), chunkNum)

		resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(buf[:n]))
		if err != nil {
			return err
		}
		resp.Body.Close()

		fmt.Printf("Uploaded chunk %d for %s\n", chunkNum, filePath)
		chunkNum++
	}

	return nil
}

func isUserFolder(name string) bool {
	systemFolders := []string{
		"$RECYCLE.BIN",
		"System Volume Information",
		"Program Files",
		"Program Files (x86)",
		"Windows",
	}
	nameLower := strings.ToLower(name)
	for _, sys := range systemFolders {
		if strings.ToLower(sys) == nameLower {
			return false
		}
	}
	return true
}

func readDriveFiles(rootDir string) ([]string, error) {
	files, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	var userfiles []string
	for _, file := range files {
		if file.IsDir() && isUserFolder(file.Name()) {
			userfiles = append(userfiles, file.Name())
		}
	}
	return userfiles, nil
}

func main() {
	rootDir := "D:/"

	dirFiles, err := readDriveFiles(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, folderName := range dirFiles {
		folderPath := filepath.Join(rootDir, folderName)

		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fmt.Println("Uploading:", path)
				if err := uploadFileInChunks(path); err != nil {
					fmt.Println("Error uploading", path, ":", err)
				}
			}
			return nil
		})

		if err != nil {
			fmt.Println("Error walking folder:", folderPath, ":", err)
		}
	}

	fmt.Println("Upload completed.")
}
