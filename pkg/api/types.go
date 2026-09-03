package api

import "time"

type StatusToday struct {
	Data struct {
		GrandTotal struct {
			TotalSecs float64 `json:"total_seconds"`
			Text      string  `json:"text"`
		} `json:"grand_total"`

		Goal struct {
			TargetSecs        float64 `json:"target_seconds"`
			TrackedSecs       float64 `json:"tracked_seconds"`
			CompletionPercent float64 `json:"completion_percent"`
			Complete          bool    `json:"complete"`
		} `json:"goal"`
	} `json:"data"`
}

type StatItem struct {
	Key  string  `json:"key"`
	Total float64 `json:"total"`
}

type WakaResponse struct {
		UserID string `json:"user_id"`
		From time.Time `json:"from"`
		To time.Time `json:"to"`
		Projects []StatItem `json:"projects"`
		Languages []StatItem `json:"languages"`
}


type Last7Days struct {
	Data struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}