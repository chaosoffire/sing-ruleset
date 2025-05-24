package constants

import (
	"fmt"
	"log"
	"os"
)

var workPath string
var configFilePath string

func getWorkPath() (workPath string, err error) {
	if workPath == "" {
		workPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	return workPath, nil
}

func setWorkPath(path string) error {
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	workPath = path
	return nil
}

func getConfigFilePath() (configFilePath string, err error) {
	if configFilePath == "" {
		workPath, err := getWorkPath()
		if err != nil {
			log.Fatalln("get work path error:", err)
		}
		configFilePath = workPath + "/config.json"
	}
	return configFilePath, nil
}

func setConfigFilePath(path string) error {
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	configFilePath = path
	return nil
}
