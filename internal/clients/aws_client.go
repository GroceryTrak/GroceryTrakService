package clients

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	CognitoClient *cognitoidentityprovider.Client
	UserPoolID    string
	ClientID      string
)

func InitAWSCognito() {
	var err error

	// Load AWS config
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Initialize Cognito client
	CognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	// Load user pool id and client id from env vars (or your config system)
	UserPoolID = os.Getenv("COGNITO_USER_POOL_ID")
	ClientID = os.Getenv("COGNITO_CLIENT_ID")

	if UserPoolID == "" || ClientID == "" {
		log.Fatalf("COGNITO_USER_POOL_ID or COGNITO_CLIENT_ID environment variables not set")
	}
}
