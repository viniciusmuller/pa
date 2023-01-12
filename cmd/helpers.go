package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/adrg/xdg"
)

var dataDirectory = path.Join(xdg.DataHome, "/pa")
var dataFilePath = path.Join(dataDirectory, "data.json")

func readDataFile() (PaData, error) {
	if _, err := os.Stat(dataDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dataDirectory, os.ModePerm)
		if err != nil {
			return PaData{}, err
		}
	}

	content, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return PaData{}, fmt.Errorf("error opening data file: %w", err)
	}

	var data PaData
	err = json.Unmarshal(content, &data)
	if err != nil {
		return PaData{}, fmt.Errorf("error parsing data file: %w", err)
	}

	return data, nil
}

func writeDataFile(data PaData) error {
	file, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(dataFilePath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
