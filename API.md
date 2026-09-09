# Hackatime API client

This package wraps the Hackatime/WakaTime-compatible API used by the TUI. It reads local WakaTime credentials, builds an authenticated HTTP client, and exposes small helper methods for the endpoints the app depends on.

## Configuration

The config loader reads your local WakaTime config file at `~/.wakatime.cfg` and extracts the API key from the `[settings]` section.

```go
creds, err := api.LoadWakatimeConfig()
if err != nil {
    panic(err)
}

client := api.NewClient(creds)
```

`Creds` is defined as:

```go
type Creds struct {
    APIKey string
    APIURL string
}
```

`LoadWakatimeConfig()` sets:

- `APIKey` from `settings.api_key`
- `APIURL` to `https://hackatime.hackclub.com/api/hackatime/v1`

## Client

```go
type Client struct {
    creds      *Creds
    httpClient *http.Client
}
```

Use `NewClient(creds *Creds)` to create the API wrapper.

### Request methods

#### `GetStatusToday()`

Fetches the current statusbar summary for today.

```go
status, err := client.GetStatusToday()
if err != nil {
    return err
}

fmt.Println(status.Data.GrandTotal.Text)
```

Endpoint:

- `GET {APIURL}/users/current/statusbar/today`
- Authorization: `Bearer {APIKey}`

Response structure:

```go
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
```

#### `GetLast7Days()`

Fetches the last 7 days aggregate stats object.

```go
stats, err := client.GetLast7Days()
if err != nil {
    return err
}

fmt.Println(stats.Data.TotalSecs)
fmt.Println(stats.Data.HumanReadableTotal)
```

Endpoint:

- `GET {APIURL}/users/current/stats/last_7_days`

Response type:

```go
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
```

`Days7Thing` represents one ranked item in the stats payload:

```go
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
```

#### `GetWakaStats(interval StatsInterval)`

This helper fetches a summary from the Hackatime summary endpoint. It first calls `GetLast7Days()` to obtain the current user ID, then requests:

- `https://hackatime.hackclub.com/api/summary?user_id={UserID}&interval={interval}`

```go
summary, err := client.GetWakaStats(api.StatsInterval_AllTime)
if err != nil {
    return err
}

fmt.Println(summary.UserID)
fmt.Println(summary.Projects[0].Key, summary.Projects[0].Total)
```

## Stats intervals

`StatsInterval` is a string alias with these supported values:

```go
type StatsInterval string

const (
    StatsInterval_Today        StatsInterval = "today"
    StatsInterval_Yesterday    StatsInterval = "yesterday"
    StatsInterval_Week         StatsInterval = "week"
    StatsInterval_7Days        StatsInterval = "7_days"
    StatsInterval_Last7Days    StatsInterval = "last_7_days"
    StatsInterval_Month        StatsInterval = "month"
    StatsInterval_30Days       StatsInterval = "30_days"
    StatsInterval_Last30Days   StatsInterval = "last_30_days"
    StatsInterval_6Months      StatsInterval = "6_months"
    StatsInterval_Last6Months  StatsInterval = "last_6_months"
    StatsInterval_Year         StatsInterval = "year"
    StatsInterval_12Months     StatsInterval = "12_months"
    StatsInterval_Last12Months StatsInterval = "last_12_months"
    StatsInterval_LastYear     StatsInterval = "last_year"
    StatsInterval_AllTime      StatsInterval = "all_time"
)
```

The summary API expects a valid interval string. If you pass an empty value, the code is effectively using the same pattern as the app's summary requests and should be treated as a custom case rather than a default in the client itself.

## Summary response model

```go
type StatItem struct {
    Key   string  `json:"key"`
    Total float64 `json:"total"`
}
```

```go
type WakaResponse struct {
    UserID    string     `json:"user_id"`
    From      time.Time  `json:"from"`
    To        time.Time  `json:"to"`
    Projects  []StatItem `json:"projects"`
    Languages []StatItem `json:"languages"`
}
```

This is a compact summary format:

- `UserID`: current user identifier
- `From` / `To`: window start/end timestamps
- `Projects`: ranked project totals
- `Languages`: ranked language totals

Example values:

```json
{
  "user_id": "12345",
  "from": "2024-01-01T00:00:00Z",
  "to": "2024-01-07T00:00:00Z",
  "projects": [
    {"key": "hackatime-tui", "total": 43200}
  ],
  "languages": [
    {"key": "Go", "total": 30000}
  ]
}
```

## Full example

```go
package main

import (
    "fmt"
    "github.com/three-thirds/hackatime-tui/pkg/api"
)

func main() {
    creds, err := api.LoadWakatimeConfig()
    if err != nil {
        panic(err)
    }

    client := api.NewClient(creds)

    today, err := client.GetStatusToday()
    if err != nil {
        panic(err)
    }

    summary, err := client.GetWakaStats(api.StatsInterval_7Days)
    if err != nil {
        panic(err)
    }

    fmt.Println("Today:", today.Data.GrandTotal.Text)
    fmt.Println("Top project:", summary.Projects[0].Key)
}
```

This covers the API pieces that are implemented in the client and request models and is the main contract the app relies on for status and summary data.
