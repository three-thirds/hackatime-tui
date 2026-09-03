package api

import "testing"

func Test67(t *testing.T) {
	creds, err := LoadWakatimeConfig()
	if err != nil {
		t.Fatalf("Failed to load Wakatime config: %v", err)
	}

	client := NewClient(creds)
	bleh, err := client.GetLast7Days()

	t.Logf("bleh %v", bleh.Data.UserID)
}	