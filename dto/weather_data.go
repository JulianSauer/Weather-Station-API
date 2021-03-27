package dto

type WeatherData struct {
    MessageId     string  `json:"messageId"`
    Timestamp     string  `json:"timestamp"`
    Temperature   float64 `json:"temperature"`
    Humidity      uint8   `json:"humidity"`
    WindSpeed     float64 `json:"windSpeed"`
    GustSpeed     float64 `json:"gustSpeed"`
    Rain          float64 `json:"rain"`
    WindDirection float32 `json:"windDirection"`
}
