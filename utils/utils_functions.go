package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadJSONFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func IsFileEmpty(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fileInfo.Size() == 0, nil
}

func WriteJSONFile(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to update new hotel data")
	}
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	return err
}

func SliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func ConvertInterfaceToString(data interface{}) string {
	return strings.TrimSpace(fmt.Sprintf("%v", data))
}
