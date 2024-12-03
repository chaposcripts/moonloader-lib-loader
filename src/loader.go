package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func isFileInSelectedList(file string, selectedModules []string) (bool, string) {
	for _, listedModule := range selectedModules {
		if strings.HasPrefix(file, listedModule+"/") {
			return true, listedModule
		}
	}
	return false, ""
}

func extractAllFiles(selectedLibs []string, zipPath, dest string) ([]string, error) {
	installedFiles := []string{}
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return installedFiles, err
	}
	for _, file := range zipFile.File {
		isSelected, moduleName := isFileInSelectedList(file.Name, selectedLibs)
		if !isSelected {
			fmt.Println("File not selected, skipping", file.Name)
			continue
		}
		filePath := fmt.Sprintf("%s\\lib\\%s", dest, strings.Replace(file.Name, moduleName+"/", "", 1))
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
				return installedFiles, err
			}
			bytes, err := io.ReadAll(reader)
			if err != nil {
				return installedFiles, err
			}

			err = os.WriteFile(filePath, bytes, 0644)
			if err != nil {
				return installedFiles, err
			}
			fmt.Println("File", filePath, "created")
			installedFiles = append(installedFiles, filePath)
		}
	}
	return installedFiles, nil
}
