package middlewares

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/controllers"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
			return []byte(config.JWT_SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			controllers.ReturnUnauthorized(c, "jwt token 过期无效")
			c.Abort()
			return
		}

		// 将 Token 信息存储到上下文中，便于后续处理
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("jwt_payload", claims["payload"])
		} else {
			controllers.ReturnServerError(c, "存储 jwt token 到上下文失败")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AccountValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtPayload, _ := c.Get("jwt_payload")
		payloadMap, ok := jwtPayload.(map[string]interface{})
		if !ok {
			fmt.Println("jwtPayload 不是一个 map[string]interface{} 类型")
			c.JSON(403, gin.H{
				"error":   "FORBIDDEN",
				"message": "访问禁止",
			})
			c.Abort()
			return
		}
		fmt.Println(payloadMap["account"])
		contentType := c.GetHeader("Content-Type")

		if contentType == "application/json" {
			var jsonData map[string]interface{}
			if err := c.ShouldBindJSON(&jsonData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			if jsonData["account"] != payloadMap["account"] {
				fmt.Println(jsonData["account"], " ", payloadMap["account"])
				c.JSON(403, gin.H{
					"error":   "FORBIDDEN",
					"message": "访问禁止",
				})
				c.Abort()
				return
			}
			c.Set("jsondata", jsonData)

		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			formData, err := c.MultipartForm()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			// fmt.Println(formData)
			account := formData.Value["account"]
			// fmt.Println(account[0])
			if account[0] != payloadMap["account"] {
				c.JSON(403, gin.H{
					"error":   "FORBIDDEN",
					"message": "访问禁止",
				})
				c.Abort()
				return
			}
			c.Set("multipartForm", formData)
		} else {
			fmt.Println("maybe do not have a body")
			account := c.Query("account")
			if account != payloadMap["account"] {
				c.JSON(403, gin.H{
					"error":   "FORBIDDEN",
					"message": "访问禁止",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
