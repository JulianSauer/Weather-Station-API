package secrets

import (
    "encoding/json"
    "github.com/JulianSauer/Weather-Station-API/dto"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

var secretId = "WeatherStation"
var region = "eu-central-1"

func Get() (*dto.Secrets, error) {
    secretsManager := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
    secretValue, e := secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretId})
    if e != nil {
        return nil, e
    }
    secretsAsJson := []byte(*secretValue.SecretString)
    secrets := dto.Secrets{}
    e = json.Unmarshal(secretsAsJson, &secrets)
    if e != nil {
        return nil, e
    }
    return &secrets, nil
}
