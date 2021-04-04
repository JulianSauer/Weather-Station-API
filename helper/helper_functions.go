package helper

import (
    "encoding/json"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "net/http"
)

func ServerError(e error) (events.APIGatewayProxyResponse, error) {
    fmt.Println(e.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func OkResponse(response interface{}) (events.APIGatewayProxyResponse, error) {
    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return ServerError(e)
    }
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers: map[string]string{
            "Content-Type":                "application/json",
            "Access-Control-Allow-Origin": "*",
        },
        Body:       string(jsonResponse),
    }, nil
}
