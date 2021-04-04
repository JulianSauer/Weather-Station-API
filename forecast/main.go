package main

import (
    "errors"
    "github.com/JulianSauer/Weather-Station-API/clients/tomorrowio"
    "github.com/JulianSauer/Weather-Station-API/helper"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
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
            if result, e := tomorrowio.Next24Hours(); e != nil {
                return helper.ServerError(e)
            } else {
                return helper.OkResponse(result)
            }
        case "daily":
            if result, e := tomorrowio.Next5Days(); e != nil {
                return helper.ServerError(e)
            } else {
                return helper.OkResponse(result)
            }
        }
    }
    return helper.ServerError(errors.New("provider or resolution not supported"))
}
