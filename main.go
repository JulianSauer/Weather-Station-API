package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "net/http"
    "time"
)

const TABLE = "WeatherStation"

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

    location, _ := time.LoadLocation("Europe/Berlin")
    t := time.Now().In(location)
    t = t.Add(-29 * time.Minute) // Sensors are updated every 15 minutes
    timestamp := t.Format("20060102-150405")

    result, e := queryLatest(timestamp)
    if e != nil {
        return serverError(e)
    }

    response := dto.WeatherData{}
    e = dynamodbattribute.UnmarshalMap(result, &response)
    if e != nil {
        return serverError(e)
    }
    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return serverError(e)
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

    result, e := queryFiltered(beginDate, endDate)
    response := make([]dto.WeatherData, len(result))
    for i, item := range result {
        data := dto.WeatherData{}
        e = dynamodbattribute.UnmarshalMap(item, &data)
        if e != nil {
            return serverError(e)
        }
        response[i] = data
    }

    jsonResponse, e := json.Marshal(response)
    if e != nil {
        return serverError(e)
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

func serverError(e error) (events.APIGatewayProxyResponse, error) {
    fmt.Println(e.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func dbConnection() *dynamodb.DynamoDB {
    session := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    return dynamodb.New(session)
}

func queryLatest(timestamp string) (map[string]*dynamodb.AttributeValue, error) {
    db := dbConnection()

    query := &dynamodb.QueryInput{
        TableName: aws.String(TABLE),
        KeyConditions: map[string]*dynamodb.Condition{
            "source": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String("WeatherStation"),
                    },
                },
            },
            "timestamp": {
                ComparisonOperator: aws.String("GT"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String(timestamp),
                    },
                },
            },
        },
    }
    result, e := db.Query(query)
    if e != nil {
        return nil, e
    }
    size := len(result.Items)
    if size > 1 {
        fmt.Printf("Found more results than expected: %d\n", size)
        fmt.Println("Picking any")
    } else if size < 1 {
        return nil, errors.New("didn't find any results")
    }
    return result.Items[0], nil
}

func queryFiltered(beginDate string, endDate string) ([]map[string]*dynamodb.AttributeValue, error) {
    db := dbConnection()

    query := &dynamodb.QueryInput{
        TableName: aws.String(TABLE),
        KeyConditions: map[string]*dynamodb.Condition{
            "source": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String("WeatherStation"),
                    },
                },
            },
            "timestamp": {
                ComparisonOperator: aws.String("BETWEEN"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String(endDate),
                    },
                    {
                        S: aws.String(beginDate),
                    },
                },
            },
        },
    }

    result, e := db.Query(query)
    if e != nil {
        return nil, e
    }
    fmt.Printf("Found %d items", len(result.Items))
    return result.Items, nil
}
