package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Record 记录请求错误
func Record(c *gin.Context, msg, err string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":  400,
		"msg":   msg,
		"error": err,
	})
}

// Response 正常返回请求结果
func Response(c *gin.Context, msg, resp string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 200,
		"msg":  msg,
		"resp": "resppppp",
	})
}
