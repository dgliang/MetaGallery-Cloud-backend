package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonStruct struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func ReturnSuccess(c *gin.Context, status string, msg string, data ...interface{}) {
	json := &JsonStruct{Code: http.StatusOK, Status: status, Msg: msg}
	if data != nil {
		json.Data = data
	}

	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, status string, msg string) {
	json := &JsonStruct{Code: http.StatusOK, Status: status, Msg: msg}

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
