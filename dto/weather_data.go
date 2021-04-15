package dto

type WeatherData struct {
    Source        string   `json:"source"`
    Timestamp     string   `json:"timestamp"`
    Temperature   []string `json:"temperature,omitempty"`
    Humidity      []string `json:"humidity,omitempty"`
    WindSpeed     []string `json:"windSpeed,omitempty"`
    GustSpeed     []string `json:"gustSpeed,omitempty"`
    Rain          []string `json:"rain,omitempty"`
    WindDirection []string `json:"windDirection,omitempty"`
    DataFor       []string `json:"dataFor,omitempty"`
}
