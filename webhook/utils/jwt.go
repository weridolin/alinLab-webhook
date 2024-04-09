package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenToken(uuid string, key string) string {
	jwt_token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		// "id":           uuid,
		"websocket_id": uuid,
		"exp":          time.Now().Add(time.Hour * 24 * 1).Unix(),
		"app":          "site.alinlab.webhook",
	}
	// Sign and get the complete encoded token as a string
	token, err := jwt_token.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err, key)
	}
	return token
}
