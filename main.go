package initializers

import (
	"database/sql"
	"github.com/PoolahGit/poolah-api-lib/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InitAWS(userPoolId string, appClientId string, appClientSecret string) *config.AwsConfig {

	conf := &aws.Config{Region: aws.String("us-east-1")}
	mySession := session.Must(session.NewSession(conf))

	a := &config.AwsConfig{
		CognitoClient:   cognito.New(mySession),
		UserPoolID:      userPoolId,
		AppClientID:     appClientId,
		AppClientSecret: appClientSecret,
	}

	return a

}

func InitDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to PlanetScale!")
	return db
}
