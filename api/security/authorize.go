package security

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Email string `json:"email" bson:"email"`
	jwt.StandardClaims
}

func ValidateToken(token string) *UserClaims {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("JWT_SECRET"), nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims
	}

	return nil
}

func Authorize(ctx *gin.Context) {
	token, err := ctx.Cookie("token")

	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})

		return
	}

	claims := ValidateToken(token)

	if claims == nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})

		return
	}

	ctx.Set("email", claims.Email)

	ctx.Next()
}
