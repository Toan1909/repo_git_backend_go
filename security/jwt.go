package security

import (
	"mygomod/model"
	"time"
	"github.com/dgrijalva/jwt-go"
)
const SECRET_KEY ="sdjvnsdfjsdofbsdfidddsdsdssd"
func GenToken(user model.User) (string,error){
	claims:= &model.JwtCustomclaims{
		UserId: user.UserId,
		Role:user.Role,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
		},
	}
	token :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	result,err:=token.SignedString([]byte(SECRET_KEY))
	if err!=nil {
		return "",err
	}
	return result,nil
}