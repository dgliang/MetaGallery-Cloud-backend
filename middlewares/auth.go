package middlewares

import (
	"MetaGallery-Cloud-backend/controllers"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secretKey = ""

func init() {
	godotenv.Load()
	secretKey = os.Getenv("JWT_SECRET_KEY")
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			controllers.ReturnUnauthorized(c, "未提供 jwt token")
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			controllers.ReturnUnauthorized(c, "提供的 jwt token 格式错误")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			controllers.ReturnUnauthorized(c, "jwt token 过期无效")
			c.Abort()
			return
		}

		// 将 Token 信息存储到上下文中，便于后续处理
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("payload", claims["payload"])
		} else {
			controllers.ReturnServerError(c, "存储 jwt token 到上下文失败")
			c.Abort()
			return
		}

		c.Next()
	}
}
