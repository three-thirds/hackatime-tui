package api

import "testing"

func Test41(t *testing.T) {
	creds, err := LoadWakatimeConfig()
	if err != nil {
		t.Fatalf("Failed to load Wakatime config: %v", err)
	}

	client := NewClient(creds)
	bleh, err := client.GetWakaStats(StatsInterval_AllTime)

	t.Logf("bleh %v", bleh.From)
	t.Logf("bleg: %s (%v seconds)", bleh.Projects[0].Key, bleh.Projects[0].Total)
}	