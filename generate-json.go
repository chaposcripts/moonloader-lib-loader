package main

import (
	"archive/zip"
	"encoding/json"
	"os"
	"fmt"
	"strings"
)

func main() {
	modules := []string{}
	outFile := "list.json"
	zipFile, err := zip.OpenReader("libs.zip")
	if err != nil {
		panic(err)
	}
	for _, file := range zipFile.File {
		if file.FileInfo().IsDir() {
			splitted := strings.Split(file.Name, "/")
			if len(splitted) == 2 {
				modules = append(modules, splitted[0])
				fmt.Println("MODULE", splitted[0])
			}
		}
	}
	bytes, err := json.Marshal(modules)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(outFile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
