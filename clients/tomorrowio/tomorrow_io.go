package tomorrowio

import (
    "encoding/json"
    "errors"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/JulianSauer/Weather-Station-API/secrets"
    "github.com/go-resty/resty/v2"
    "net/http"
    "time"
)

const ISO8601 = "2006-01-02T15:04:05Z"

const HOURLY = "&timesteps=1h"
const DAILY = "&timesteps=1d"

var API = "https://api.tomorrow.io/v4/timelines?" +
    "units=metric" +
    "&timezone=Europe/Berlin" +
    "&fields=temperature" +
    "&fields=precipitationProbability"

type Forecast struct {
    Data struct {
        Timelines []struct {
            Timestep  string    `json:"timestep"`
            StartTime time.Time `json:"startTime"`
            EndTime   time.Time `json:"endTime"`
            Intervals []struct {
                StartTime time.Time `json:"startTime"`
                Values    struct {
                    Temperature              float64 `json:"temperature"`
                    PrecipitationProbability float64 `json:"precipitationProbability"`
                } `json:"values"`
            } `json:"intervals"`
        } `json:"timelines"`
    } `json:"data"`
}

func Next24Hours() ([]dto.ForecastResult, error) {
    now := time.Now()
    location, _ := time.LoadLocation("Europe/Berlin")
    startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, location).In(time.UTC)
    endTime := startTime.AddDate(0, 0, 1)
    return queryTomorrowIOAPI(HOURLY, startTime.Format(ISO8601), endTime.Format(ISO8601))
}

func Next5Days() ([]dto.ForecastResult, error) {
    startTime := time.Now().In(time.UTC).Format(ISO8601)
    endTime := time.Now().In(time.UTC).AddDate(0, 0, 5).Format(ISO8601)
    return queryTomorrowIOAPI(DAILY, startTime, endTime)
}

func queryTomorrowIOAPI(timeInterval string, startTime string, endTime string) ([]dto.ForecastResult, error) {
    client := resty.New()
    baseUrl, e := createBaseUrl()
    if e != nil {
        return nil, e
    }
    url := *baseUrl +
        timeInterval +
        "&startTime=" + startTime +
        "&endTime=" + endTime

    response, e := client.R().Get(url)
    if e != nil {
        return nil, e
    }
    if response.StatusCode() != http.StatusOK {
        return nil, errors.New(string(response.Body()))
    }

    forecast := Forecast{}
    if e := json.Unmarshal(response.Body(), &forecast); e != nil {
        return nil, e
    }
    forecastData := forecast.Data.Timelines[0].Intervals
    result := make([]dto.ForecastResult, len(forecastData))
    for i, entry := range forecastData {
        result[i] = dto.ForecastResult{
            Timestamp:                entry.StartTime.Format("20060102-150405"),
            Temperature:              entry.Values.Temperature,
            PrecipitationProbability: entry.Values.PrecipitationProbability,
        }
    }
    return result, nil
}

func createBaseUrl() (*string, error) {
    s, e := secrets.Get()
    if e != nil {
        return nil, e
    }
    url := API +
        "&location=" + s.Latitude + "," + s.Longitude +
        "&apikey=" + s.ApiKeyTomorrowIO
    return &url, nil
}
