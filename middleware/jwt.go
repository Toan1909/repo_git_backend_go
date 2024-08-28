package middleware

import (
	"mygomod/model"
	"mygomod/security"
	jwt"github.com/golang-jwt/jwt/v4"
	jwtMiddleware "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func JWTMiddleWare() echo.MiddlewareFunc{
		config:= jwtMiddleware.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims{
				return new(model.JwtCustomclaims)
			},
			SigningKey: []byte(security.SECRET_KEY),
			//contextKey mặc định là "user",có thể truyền string vào contextKey để thay đối 
		}
		return jwtMiddleware.WithConfig(config)
}
