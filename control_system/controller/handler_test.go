package controller

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestJwtSign(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("hello"))
	assert.Nil(t, err)
	verToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("hello"), nil
	})
	assert.Nil(t, err)
	assert.True(t, verToken.Valid)
}
