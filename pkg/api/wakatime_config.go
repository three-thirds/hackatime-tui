package api

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Creds struct {
	APIKey string
	APIURL string
}

func LoadWakatimeConfig() (*Creds, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".wakatime.cfg")
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to blehhhhh")
	}

	section := cfg.Section("settings")
	apiKey := section.Key("api_key").String()
	if apiKey == "" {
		return nil, fmt.Errorf("bleh2")
	}

	apiURL := "https://hackatime.hackclub.com/api/hackatime/v1"

	return &Creds{APIKey: apiKey, APIURL: apiURL}, nil
}
