package util

import (
	"github.com/gin-gonic/gin"
)

func GetAccessToken(c *gin.Context) (string, error) {
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header is required")
	// }
	// tokenString := authHeader[len("Bearer "):]

	// if tokenString == "" {
	// 	return "", errors.New("Bearer token is required")
	// }
	tokenString, err := c.Cookie("access_token")
	// fmt.Println("Cookie access_token:", tokenString)
	if err != nil {

		return "", err
	}

	return tokenString, nil

}
