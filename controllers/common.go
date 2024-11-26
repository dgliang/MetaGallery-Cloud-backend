package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonStruct struct {
	Code   int         `json:"code"`   // 状态码如 200，400，401
	Status string      `json:"status"` // 状态如 SUCCESS，FAILED
	Msg    string      `json:"msg"`    // 状态信息
	Data   interface{} `json:"data"`   // 数据
}

func ReturnSuccess(c *gin.Context, status string, msg string, data ...interface{}) {
	json := &JsonStruct{Code: http.StatusOK, Status: status, Msg: msg}
	if len(data) > 0 {
		json.Data = data[0]
	}

	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, status string, msg string) {
	json := &JsonStruct{Code: http.StatusBadRequest, Status: status, Msg: msg}

	c.JSON(http.StatusOK, json)
}

func ReturnServerError(c *gin.Context, msg string) {
	json := &JsonStruct{Code: http.StatusInternalServerError, Msg: msg}

	c.JSON(http.StatusInternalServerError, json)
}

func ReturnUnauthorized(c *gin.Context, msg string) {
	json := &JsonStruct{Code: http.StatusUnauthorized, Msg: msg}

	c.JSON(http.StatusUnauthorized, json)
}
