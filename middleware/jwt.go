package middleware

import (
	"net/http"
	"time"
	"user-system/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("你的Secret")

func GenerateToken(Name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": Name,
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(Secret)
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.GetHeader("Authorization")
		_, err := config.RDB.Get(
			config.Ctx,
			tokenstring,
		).Result()
		if err != nil {
			c.JSON(401, gin.H{
				"error": "token失效",
			})
			c.Abort()
		}
		if tokenstring == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "未登录",
			})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "token无效",
			})
			c.Abort()
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "claim断言失败",
			})
		}
		username := claim["name"]
		c.Set("username", username)
		c.Next()
	}
}
