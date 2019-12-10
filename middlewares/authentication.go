package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/utils"
	"net/http"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		if bearer := ctx.GetHeader("Authorization"); len(bearer) > 7 {
			token = bearer[7:]
		}

		claims := jwt.StandardClaims{}

		if _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("environment-var"), nil
		}); err != nil {
			utils.SetResponse(ctx, http.StatusUnauthorized, "invalid token")

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
