package model

import (
		"github.com/dgrijalva/jwt-go"
)


type JwtCustomclaims struct {
	UserId string
	Role   string
	jwt.StandardClaims
}