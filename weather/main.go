package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-API/db"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/JulianSauer/Weather-Station-API/helper"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "net/http"
    "strconv"
)

func main() {
    lambda.Start(router)
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

func getLatestWeatherData() (events.APIGatewayProxyResponse, error) {
    fmt.Println("Looking up latest data")

    result, e := db.QueryLatest("WeatherStation")
    if e != nil {
        return helper.ServerError(e)
    }

    response := dto.WeatherData{}
    e = dynamodbattribute.UnmarshalMap(result, &response)
    if e != nil {
        return helper.ServerError(e)
    }
    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return helper.ServerError(e)
    }
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers: map[string]string{
            "Content-Type":                "application/json",
            "Access-Control-Allow-Origin": "*",
        },
        Body: string(jsonResponse),
    }, nil
}

func getWeatherDataFiltered(beginDate string, endDate string) (events.APIGatewayProxyResponse, error) {
    fmt.Printf("Filtering: %s - %s\n", beginDate, endDate)

    result, e := db.QueryFiltered(beginDate, endDate, "WeatherStation")
    response := make([]dto.WeatherData, len(result))
    var firstRainData float64
    for i, item := range result {
        data := dto.WeatherData{}
        e = dynamodbattribute.UnmarshalMap(item, &data)
        if e != nil {
            return helper.ServerError(e)
        }
        firstRainData, data.Rain[0], e = zeroRainData(i, firstRainData, data.Rain[0])
        if e != nil {
            return helper.ServerError(e)
        }
        response[i] = data
    }

    return helper.OkResponse(response)
}

func zeroRainData(i int, firstRainData float64, currentRainData string) (float64, string, error) {
    if i == 0 {
        if firstRain, e := strconv.ParseFloat(currentRainData, 64); e != nil {
            return firstRain, "", e
        } else {
            return firstRain, "0.0", nil
        }
    } else {
        if currentRain, e := strconv.ParseFloat(currentRainData, 64); e != nil {
            return firstRainData, "", e
        } else {
            currentRain -= firstRainData
            return firstRainData, fmt.Sprintf("%f", currentRain), nil
        }
    }
}
