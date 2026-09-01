package api

import "testing"

func TestGoog(t *testing.T) {
	creds, err := LoadWakatimeConfig()
	if err != nil {
		return
	}
	client := NewClient(creds)
	status, err := client.GetStatusToday()

	t.Logf("bleh %v", status.Data.GrandTotal.TotalSecs)
}
