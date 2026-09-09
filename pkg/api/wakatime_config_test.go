package api

import "testing"

func TestBleh(t *testing.T) {
	creds, err := LoadWakatimeConfig()
	if err != nil {
		t.Fatalf("Failed to load Wakatime config: %v", err)
	}

	t.Logf("API Key: %s", creds.APIKey)
	t.Logf("API URL: %s", creds.APIURL)
}
