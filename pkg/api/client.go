package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	creds      *Creds
	httpClient *http.Client
}

func NewClient(creds *Creds) *Client {
	return &Client{
		creds:      creds,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetStatusToday() (*StatusToday, error) {
	url := c.creds.APIURL + "/users/current/statusbar/today"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.creds.APIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bleh bleh bleh skill issue")
	}

	var status StatusToday
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) GetLast7Days() (*Last7Days, error) {
	return c.GetStats("last_7_days")
}

// GetStats fetches aggregate stats for a WakaTime-compatible range
// (e.g. "today", "last_7_days", "last_30_days", "all_time").
func (c *Client) GetStats(rangeName string) (*Last7Days, error) {
	if rangeName == "" {
		rangeName = "last_7_days"
	}
	url := c.creds.APIURL + "/users/current/stats/" + rangeName

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.creds.APIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get %s stats: %s", rangeName, response.Status)
	}

	var status Last7Days
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

type StatsInterval string
const (
	StatsInterval_Today StatsInterval = "today"
	StatsInterval_Yesterday StatsInterval = "yesterday"
	StatsInterval_Week StatsInterval = "week"
	StatsInterval_7Days StatsInterval = "7_days"
	StatsInterval_Last7Days StatsInterval = "last_7_days"
	StatsInterval_Month StatsInterval = "month"
	StatsInterval_30Days StatsInterval = "30_days"
	StatsInterval_Last30Days StatsInterval = "last_30_days"
	StatsInterval_6Months StatsInterval = "6_months"
	StatsInterval_Last6Months StatsInterval = "last_6_months"
	StatsInterval_Year StatsInterval = "year"
	StatsInterval_12Months StatsInterval = "12_months"
	StatsInterval_Last12Months StatsInterval = "last_12_months"
	StatsInterval_LastYear StatsInterval = "last_year"
	StatsInterval_AllTime StatsInterval = "all_time"
)

func (c *Client) GetWakaStats(interval StatsInterval) (*WakaResponse, error) {
	last7, err := c.GetLast7Days()
	if err != nil {
		return nil, err
	}
	UserID := last7.Data.UserID

	url := "https://hackatime.hackclub.com/api/summary?user_id=" + UserID + "&interval=" + string(interval)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get waka stats: %s", response.Status)
	}

	var wakaStats WakaResponse
	if err := json.NewDecoder(response.Body).Decode(&wakaStats); err != nil {
		return nil, err
	}

	return &wakaStats, nil
}