package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const jsonURL string = "https://raw.githubusercontent.com/chaposcripts/moonloader-lib-loader/refs/heads/main/list.json"

type LibsData map[string][]string

var dataReceived = false
var libsData LibsData

func loadData() error {
	response, err := http.Get(jsonURL)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP_STATUS_%d", response.StatusCode)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &libsData)
	if err == nil {
		dataReceived = true
	}
	return err
}

func moveSelectedLibs(selectedLibs []string) error {
	for libName, files := range libsData {
		for _, it := range selectedLibs {
			if it == libName {
				for _, fname := range files {
					defaultFile := fmt.Sprintf("./temp_libs_zip/%s", fname)
					newFile := fmt.Sprintf("./moonloader/lib/%s", fname)
					// fmt.Println("MOVING \"%s\" to \"%s\"", fmt.Sprintf("./temp_libs_zip/%s", fname), fmt.Sprintf("./moonloader/lib/%s", fname))
					isDir := isDirectory(defaultFile)
					if isDir {
						copyDir(defaultFile, newFile)
					} else {
						copyFile(defaultFile, newFile)
					}
				}
			}
		}
	}
	return nil
}
