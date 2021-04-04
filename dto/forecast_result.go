package dto

type ForecastResult struct {
    Timestamp   string  `json:"timestamp"`
    Temperature float64 `json:"temperature"`
    PrecipitationProbability int     `json:"precipitationProbability"`
}
