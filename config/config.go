package config

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type AwsConfig struct {
	CognitoClient   *cognito.CognitoIdentityProvider
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
	Token           string
}
