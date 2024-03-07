package app

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func LoadConfiguration() (Config, error) {
	var config Config

	body, err := os.ReadFile("./files/config.json")
	if err != nil {
		return config, errors.New("file config not found")
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func LoadFile(e Exploit) ([]string, error) {
	body, err := os.ReadFile(e.FilePath)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(strings.ReplaceAll(string(body), "\r", ""), "\n")
	return paths, nil
}

func WriteResult(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return nil
}

func FilterUrl(url string) string {
	url = strings.ReplaceAll(url, "\r", "")
	url = strings.TrimSuffix(url, "/")
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}
