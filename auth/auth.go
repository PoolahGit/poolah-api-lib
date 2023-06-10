package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"

	Config "github.com/PoolahGit/poolah-api-lib/config"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type ParsedJwt struct {
	Username string
	Email    string
	Phone    string
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
func VerifyJWT(authHeaderArray []string, a *Config.AwsConfig, c *gin.Context) (*ParsedJwt, error) {

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

	username, _ := token.Get("cognito:username")
	email, _ := token.Get("email")
	phone, _ := token.Get("phone_number")

	parsedUser := &ParsedJwt{
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
