package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-API/clients/tomorrowio"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/JulianSauer/Weather-Station-API/secrets"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
    "time"
)

var location *time.Location
var topic string

func main() {
    lambda.Start(handler)
}

func handler() {
    location, _ = time.LoadLocation("Europe/Berlin")
    s, e := secrets.Get()
    if e != nil {
        fmt.Println(e.Error())
        return
    }
    topic = s.SnsDataTopic

    data, errors := getData()
    for _, e := range errors {
        fmt.Println(e.Error())
    }
    PublishSensorData(&data)
}

func getData() ([]string, []error) {
    var results []string
    var errors []error

    results, errors = formatData(results, errors, "TomorrowIO-Hourly", tomorrowio.Next24Hours)
    results, errors = formatData(results, errors, "TomorrowIO-Daily", tomorrowio.Next5Days)

    return results, errors
}

func formatData(results []string, errors []error, provider string, providerFunction func() ([]dto.ForecastResult, error)) ([]string, []error) {
    forecastResults, e := providerFunction()
    if e != nil {
        errors = append(errors, e)
    } else {

        timestamp := time.Now().In(location).Format("20060102-150405")
        var temperatures []string
        var rain []string
        var forecastTimestamps []string
        for _, forecastResult := range forecastResults {
            forecastTimestamps = append(forecastTimestamps, forecastResult.Timestamp)
            temperatures = append(temperatures, fmt.Sprintf("%.1f", forecastResult.Temperature))
            rain = append(rain, fmt.Sprintf("%.0f", forecastResult.PrecipitationProbability))
        }

        minedWeatherForecast := dto.WeatherData{
            Source:        provider,
            Timestamp:     timestamp,
            Temperature:   temperatures,
            Rain:          rain,
            DataFor:       forecastTimestamps,
            Humidity:      []string{},
            WindSpeed:     []string{},
            GustSpeed:     []string{},
            WindDirection: []string{},
        }
        resultsAsJson, e := json.Marshal(minedWeatherForecast)
        if e != nil {
            errors = append(errors, e)
        } else {
            results = append(results, string(resultsAsJson))

            latest := minedWeatherForecast
            latest.Timestamp = "latest"
            latestAsJson, _ := json.Marshal(latest)
            results = append(results, string(latestAsJson))
        }
    }
    return results, errors
}

func PublishSensorData(messages *[]string) {
    session, e := session.NewSession(&aws.Config{
        Region: aws.String("eu-central-1"),
    })

    if e != nil {
        fmt.Println(e.Error())
        return
    }

    client := sns.New(session)
    for _, message := range *messages {
        input := &sns.PublishInput{
            Message:  aws.String(message),
            TopicArn: aws.String(topic),
        }

        result, e := client.Publish(input)
        if e != nil {
            fmt.Printf("Cannot publish: %s\n", e.Error())
        } else {
            fmt.Printf("%s: %s\n", *result.MessageId, message)
        }
    }
}
