package initializers

import (
	"database/sql"
	"fmt"
	"github.com/PoolahGit/poolah-api-lib/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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

func InitDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%v", viper.Get("DSN")))

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to PlanetScale!")
	return db
}

type ForbiddenError struct {
	statusCode int
	err        string
}

type UnauthorizedError struct {
	statusCode int
	err        string
}

func (forbidden *ForbiddenError) Error() string {
	return fmt.Sprintf(
		"Status code %s, you are not authorized to call this API, or your credentials expired",
		forbidden.statusCode,
	)
}

func (unauth *UnauthorizedError) Error() string {
	return fmt.Sprintf(
		"Status code %s, you are not logged in",
		unauth.statusCode,
	)
}

// Re-usable function
func VerifyJWT(authHeaderArray []string, a *Config.AwsConfig, c *gin.Context) (*models.ParsedJwt, error) {

	if len(authHeaderArray) == 0 {
		return nil, &UnauthorizedError{
			err:        "",
			statusCode: 401,
		}
	}
	authHeader := authHeaderArray[0]

	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 {
		return nil, &UnauthorizedError{
			err:        "unauthorized",
			statusCode: 401,
		}
	}

	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, "us-east-1", a.UserPoolID)
	keySet, err := jwk.Fetch(c, formattedURL)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(
		[]byte(splitAuthHeader[1]),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)

	if err != nil {
		fmt.Println(err)
		return nil, &ForbiddenError{
			err:        "Forbidden",
			statusCode: 403,
		}
	}

	fmt.Println(token)

	username, _ := token.Get("cognito:username")
	email, _ := token.Get("email")
	phone, _ := token.Get("phone_number")

	parsedUser := &models.ParsedJwt{
		Username: fmt.Sprint(username),
		Email:    fmt.Sprint(email),
		Phone:    fmt.Sprint(phone),
	}

	return parsedUser, nil

}

func AuthHandlerFunc(a *Config.AwsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeaderSplit := c.Request.Header["Token"]
		userContext, err := VerifyJWT(authHeaderSplit, a, c)

		if err != nil {
			var forbiddenError *ForbiddenError
			var unauthorizedError *UnauthorizedError

			if errors.As(err, &forbiddenError) {
				c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
				return
			} else if errors.As(err, &unauthorizedError) {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}
		}

		c.Set("Username", userContext.Username)
		c.Set("UserEmail", userContext.Email)
		c.Set("UserPhone", userContext.Phone)

		c.Next()
	}
}
