package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "net/http"
)

func getLatestWeatherData() (events.APIGatewayProxyResponse, error) {
    fmt.Println("Looking up latest data")
    response := dto.WeatherData{
        MessageId:     "",
        Timestamp:     "",
        Temperature:   0,
        Humidity:      0,
        WindSpeed:     0,
        GustSpeed:     0,
        Rain:          0,
        WindDirection: 0,
    }

    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return serverError(e)
    }
    return events.APIGatewayProxyResponse{
        StatusCode:        http.StatusOK,
        Body:              string(jsonResponse),
    }, nil
}

func getWeatherDataFiltered(beginDate string, endDate string) (events.APIGatewayProxyResponse, error) {
    fmt.Printf("Filtering: %s - %s\n", beginDate, endDate)
    element1 := dto.WeatherData{
        MessageId:     "",
        Timestamp:     beginDate,
        Temperature:   0,
        Humidity:      0,
        WindSpeed:     0,
        GustSpeed:     0,
        Rain:          0,
        WindDirection: 0,
    }

    element2 := dto.WeatherData{
        MessageId:     "",
        Timestamp:     endDate,
        Temperature:   1,
        Humidity:      2,
        WindSpeed:     3,
        GustSpeed:     4,
        Rain:          5,
        WindDirection: 6,
    }

    response := [2]dto.WeatherData{element1, element2}
    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return serverError(e)
    }
    return events.APIGatewayProxyResponse{
        StatusCode:        http.StatusOK,
        Body:              string(jsonResponse),
    }, nil
}

func serverError(e error) (events.APIGatewayProxyResponse, error) {
    fmt.Println(e.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    beginDate, beginDateExists := request.QueryStringParameters["begin"]
    endDate, endDateExists := request.QueryStringParameters["end"]
    if beginDateExists && endDateExists {
        return getWeatherDataFiltered(beginDate, endDate)
    } else {
        return getLatestWeatherData()
    }
}

func main() {
    lambda.Start(router)
}
