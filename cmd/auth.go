package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type TokenFile struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func getToken(envName string) (string, error) {
	tokenFile, err := tokenFile(runtime.GOOS, envName)
	if err != nil {
		return "", err
	}

	raw, err := os.ReadFile(tokenFile)
	if err != nil {
		panic(err)
	}

	data := &TokenFile{}
	err = json.Unmarshal(raw, data)
	if err != nil {
		panic(err)
	}

	return data.Token, nil
}

func tokenFile(goos, envName string) (string, error) {
	switch goos {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "fermyon", fmt.Sprintf("%s.json", envName)), nil
	case "windows":
		return filepath.Join(os.Getenv("HOME"), "AppData", "Roaming", "fermyon", fmt.Sprintf("%s.json", envName)), nil
	case "linux":
		return filepath.Join(os.Getenv("HOME"), ".config", "fermyon", fmt.Sprintf("%s.json", envName)), nil
	}

	return "", fmt.Errorf("%s os not supported", goos)
}
