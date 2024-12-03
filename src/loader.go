package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
)

const zipURL string = "https://github.com/chaposcripts/moonloader-lib-loader/raw/refs/heads/main/libs.zip"

func loadZip(url, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP_STATUS_%d", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bodyBytes, 0644)
	return err
}

func isFileInList(list []string, targetFile string) bool {
	for _, moduleName := range list {
		for module, files := range libsData {
			if module == moduleName {
				for _, file := range files {
					if file == targetFile {
						return true
					}
				}
			}
		}
	}
	return false
}

func extractAllFiles(selectedLibs []string, zipPath, dest string) error {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	for _, file := range zipFile.File {
		if !isFileInList(selectedLibs, file.Name) {
			fmt.Println("File not selected, skipping", file.Name)
			continue
		}
		filePath := fmt.Sprintf("%s\\%s", dest, file.Name)
		fmt.Println("ZIP", file.Name, filePath, file.FileInfo().IsDir())
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			os.Remove(filePath)
			fmt.Println("File", filePath, "already exists, removed!")
		}
		if file.FileInfo().IsDir() {
			os.Mkdir(filePath, 0644)
			fmt.Println("Dir", filePath, "Created")
		} else {
			reader, err := file.Open()
			if err != nil {
				return err
			}
			bytes, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			err = os.WriteFile(filePath, bytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println("File", filePath, "created")
		}
	}
	return nil
}
