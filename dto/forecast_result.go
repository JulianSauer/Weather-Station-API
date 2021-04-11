package dto

type ForecastResult struct {
    Timestamp                string  `json:"timestamp"`
    Temperature              float64 `json:"temperature"`
    PrecipitationProbability float64 `json:"precipitationProbability"`
}
