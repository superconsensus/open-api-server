package middlewares

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"

	"github.com/superconsensus/open-api-server/conf"
)

func Logs() gin.HandlerFunc {

	dir := path.Base(conf.Log.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	logFilePath := conf.Log.FilePath
	logFileName := conf.Log.RouterFile
	logfile := path.Join(logFilePath, logFileName)
	configs := `{"filename":"` + logfile + `","color":true}`

	logger := logs.NewLogger(1000)
	err = logger.SetLogger(logs.AdapterFile, configs)
	if err != nil {
		panic(err)
	}
	logger.Async()

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 状态码
		statusCode := c.Writer.Status()
		if statusCode == http.StatusOK || statusCode == http.StatusNotFound {
			//成功的请求不记录
			return
		}

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Info("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
