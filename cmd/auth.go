package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Token struct {
	URL          string    `json:"url"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	Expiration   time.Time `json:"expiration"`
}

func getToken(envName string) (*Token, error) {
	tokenFile, err := getTokenFile(runtime.GOOS, envName)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(tokenFile)
	if err != nil {
		panic(err)
	}

	data := &Token{}
	err = json.Unmarshal(raw, data)
	if err != nil {
		panic(err)
	}

	return data, nil
}

func getTokenFile(goos, envName string) (string, error) {
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

func saveNewToken(envName string, token *Token) error {
	tokenFile, err := getTokenFile(runtime.GOOS, envName)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return os.WriteFile(tokenFile, raw, 0644)
}
