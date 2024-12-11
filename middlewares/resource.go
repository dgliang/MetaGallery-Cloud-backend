package middlewares

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResourceAccessAuthMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		account := c.Query("account")
		fileID := c.Query("file_id")

		// 检查用户的角色，验证其是否可以预览相应的文件
		err, avail := IsAccessValid(account, fileID)
		if err == nil && avail {
			c.Next()
		} else {
			// controllers.ReturnError(c, "FORBIDDEN", "访问被禁止")
			c.JSON(403, gin.H{
				"error":   "FORBIDDEN",
				"message": "访问禁止",
			})
			c.Abort() // 阻止请求继续执行
			return
		}
	}
}

func IsAccessValid(account, fileID string) (error, bool) {

	log.Printf("%s 希望访问资源 %s ", account, fileID)

	userID, err := models.GetUserID(account)
	if err != nil {
		return err, false
	}
	if userID == 0 {
		return fmt.Errorf("用户不存在"), false
	}

	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {

		return err, false
	}
	uintFID := uint(FID)

	Belongto := services.IsFileBelongto(userID, uintFID)
	if Belongto {
		return nil, true
	}
	return nil, false
}
