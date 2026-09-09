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
	Key   string  `json:"key"`
	Total float64 `json:"total"`
}

type WakaResponse struct {
	UserID    string     `json:"user_id"`
	From      time.Time  `json:"from"`
	To        time.Time  `json:"to"`
	Projects  []StatItem `json:"projects"`
	Languages []StatItem `json:"languages"`
}

type Days7Thing struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	DigitalTime  string  `json:"digital"`
	TextTime     string  `json:"text"`
	Hours        float64 `json:"hours"`
	Minutes      float64 `json:"minutes"`
	Seconds      float64 `json:"seconds"`
}

type Last7Days struct {
	Data struct {
		Username                  string       `json:"username"`
		UserID                    string       `json:"user_id"`
		Start                     time.Time    `json:"start"`
		End                       time.Time    `json:"end"`
		Status                    string       `json:"status"`
		TotalSecs                 float64      `json:"total_seconds"`
		DailyAverageSecs          float64      `json:"daily_average_seconds"`
		DaysIncludingHolidays     int          `json:"days_including_holidays"`
		Range                     string       `json:"range"`
		HumanReadableRange        string       `json:"human_readable_range"`
		HumanReadableTotal        string       `json:"human_readable_total"`
		HumanReadableDailyAverage string       `json:"human_readable_daily_average"`
		IsCodingActivityVisible   bool         `json:"is_coding_activity_visible"`
		IsOtherUsageVisible       bool         `json:"is_other_usage_visible"`
		Editors                   []Days7Thing `json:"editors"`
		Languages                 []Days7Thing `json:"languages"`
		Projects                  []Days7Thing `json:"projects"`
		OperatingSystems          []Days7Thing `json:"operating_systems"`
		Categories                []Days7Thing `json:"categories"`
	} `json:"data"`
}
