package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func validateToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("Incorrect JWT signing method")
	}

	// Get public key to check signature.
	// RSA_PUBLIC_KEY - environment variable responsible for JWT public key.
	encPublicKey := os.Getenv("RSA_PUBLIC_KEY")
	publicKey, err := base64.StdEncoding.DecodeString(encPublicKey)

	if err != nil {
		return nil, err
	}

	signingKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		return nil, err
	}

	return signingKey, nil
}

// AuthMiddleware is gin middleware, thats validates JWT token.
func AuthMiddleware(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")

	if authorizationHeader != "" {

		splittedHeader := strings.Split(authorizationHeader, " ")

		if len(splittedHeader) < 2 {
			c.JSON(http.StatusUnauthorized, GenericError{Error: "Invalid header content"})
			c.Abort()
			return
		}

		tokenType := splittedHeader[0]   // Assign token type.
		tokenString := splittedHeader[1] // Assign token value.

		// Validate JWT token type.
		// Token type should be equal to "Bearer".
		if tokenType != "Bearer" {
			c.JSON(http.StatusUnauthorized, GenericError{Error: "Invalid token type"})
			c.Abort()
			return
		}

		// Validate JWT token.
		token, err := jwt.Parse(tokenString, validateToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, GenericError{Error: err.Error()})
			c.Abort()
			return
		}

		if token.Valid && token.Claims.(jwt.MapClaims)["role"] == "admin" {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, GenericError{Error: "Not Authorized"})
			c.Abort()
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, GenericError{Error: "Not Authorized"})
		c.Abort()
		return
	}

}
