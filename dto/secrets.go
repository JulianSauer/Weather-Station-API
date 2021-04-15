package dto

type Secrets struct {
    ApiKeyTomorrowIO string `json:"apiKeyTomorrowIO"`
    Latitude         string `json:"latitude"`
    Longitude        string `json:"longitude"`
    SnsDataTopic     string `json:"snsDataTopic"`
}
