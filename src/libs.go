package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const jsonURL string = "https://raw.githubusercontent.com/chaposcripts/moonloader-lib-loader/refs/heads/main/list.json"

type LibsData []string

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
