package db

import (
    "errors"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

const TABLE = "WeatherStation"

func dbConnection() *dynamodb.DynamoDB {
    session := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    return dynamodb.New(session)
}

func QueryLatest(source string) (map[string]*dynamodb.AttributeValue, error) {
    db := dbConnection()

    query := &dynamodb.QueryInput{
        TableName: aws.String(TABLE),
        KeyConditions: map[string]*dynamodb.Condition{
            "source": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String(source),
                    },
                },
            },
            "timestamp": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String("latest"),
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

func QueryFiltered(beginDate string, endDate string, source string) ([]map[string]*dynamodb.AttributeValue, error) {
    db := dbConnection()

    query := &dynamodb.QueryInput{
        TableName: aws.String(TABLE),
        KeyConditions: map[string]*dynamodb.Condition{
            "source": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String(source),
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
