package main

import (
    "errors"
    "github.com/JulianSauer/Weather-Station-API/db"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/JulianSauer/Weather-Station-API/helper"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
    lambda.Start(router)
}

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    provider, providerExists := request.QueryStringParameters["from"]
    resolution, resolutionExists := request.QueryStringParameters["resolution"]
    if !resolutionExists {
        return helper.ServerError(errors.New("missing parameter: resolution"))
    }
    if !providerExists {
        return helper.ServerError(errors.New("missing parameter: provider"))
    }
    switch provider {
    case "TomorrowIO":
        switch resolution {
        case "hourly":
            result, e := db.QueryLatest("TomorrowIO-Hourly")
            if e != nil {
                return helper.ServerError(e)
            } else {
                response := dto.WeatherData{}
                e = dynamodbattribute.UnmarshalMap(result, &response)
                if e != nil {
                    return helper.ServerError(e)
                }
                return helper.OkResponse(response)
            }
        case "daily":
            result, e := db.QueryLatest("TomorrowIO-Daily")
            if e != nil {
                return helper.ServerError(e)
            } else {
                response := dto.WeatherData{}
                e = dynamodbattribute.UnmarshalMap(result, &response)
                if e != nil {
                    return helper.ServerError(e)
                }
                return helper.OkResponse(response)
            }
        }
    }
    return helper.ServerError(errors.New("provider or resolution not supported"))
}
