package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/palestine-nights/backend/src/db"
	"github.com/palestine-nights/backend/src/tools"
)

type Token struct {
	Token string `json:"token"`
}

func (server *Server) authenticateUser(c *gin.Context) {
	requestedUser := db.User{}

	if err := c.BindJSON(&requestedUser); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	user, err := db.User.FindByPassword(db.User{}, server.DB, requestedUser.UserName, requestedUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Incorrect credentials"})
		return
	}

	// Generate token
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["username"] = user.UserName
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Get signing key
	privateKey := []byte(tools.GetEnv("RSA_PRIVATE_KEY", ""))
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Get token string
	tokenString, err := token.SignedString(signingKey)
	tokenObj := Token{Token: tokenString}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, tokenObj)
	}
}

func authorized(endpoint func(c *gin.Context)) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New("incorrect jwt signing method")
				}

				// Get public key to check signature
				publicKey := []byte(tools.GetEnv("RSA_PUBLIC_KEY", ""))
				signingKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
				if err != nil {
					return nil, err
				}

				return signingKey, nil
			})

			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
			}

			if token.Valid && token.Claims.(jwt.MapClaims)["admin"] == true {
				endpoint(c)
			} else {
				c.JSON(http.StatusUnauthorized, GenericError{Error: "Not authorized"})
			}
		} else {
			c.JSON(http.StatusUnauthorized, GenericError{Error: "Not authorized"})
		}
	})
}
