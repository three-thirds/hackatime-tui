# How to use API

## Different Fetches

- Client
- Wakatime Config

### Fetches

- [Past 7 days](#past-7-days)
- [Wakatime Response](#wakatime-response)
- Status Today

## How to implement

```go
package api

func ___ () {
    creds, err := LoadWakatimeConfig()
    if err != nil {
        t.Fatalf("Failed to load Hackatime config: %v", err)
    }

    client :=NewClient(Creds)
    response, err := client.[Put fetch here]
}
```

### Data structure & Fetch params

**All Params will be listed in order of use.**

#### Wakatime Response

##### Fetch Params

- Time interval
  - StatsInterval_Today
  - StatsInterval_Yesterday
  - StatsInterval_Week
  - StatsInterval_7Days
  - StatsInterval_Last7Days
  - StatsInterval_Month
  - StatsInterval_30Days
  - StatsInterval_Last30Days
  - StatsInterval_6Months
  - StatsInterval_Last6Months
  - StatsInterval_Year
  - StatsInterval_12Months
  - StatsInterval_Last12Months
  - StatsInterval_LastYear
  - StatsInterval_AllTime

###### Note

Empty time interval will default to All Time.

##### Data response Structure

- __.UserID
- __.From
- __.To
- __.Projects[0]
  
```json
    "Key" : "Project Name"
    "Total" : "Total Seconds"
```

- __.Languages[0]

``` json
    "Key" : "Language"
    "Total" : "Total Seconds"
```

###### Ordering Note

Projects and Languages are returned from greatest to lowest according to Total
