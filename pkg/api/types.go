package api

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
