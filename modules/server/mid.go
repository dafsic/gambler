package server

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/dafsic/gambler/utils"
	"github.com/gin-gonic/gin"
)

type responseLogWriter struct {
	gin.ResponseWriter
	responseBody *bytes.Buffer
}

func (w responseLogWriter) Write(b []byte) (int, error) {
	w.responseBody.Write(b)
	return w.ResponseWriter.Write(b)
}

// Record what request and response are
func Record(l *utils.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		urlPath := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		method := ctx.Request.Method
		clientIP := ctx.ClientIP()
		ctx.Set("cip", clientIP)
		bodyBytes, _ := io.ReadAll(ctx.Request.Body)

		if raw != "" {
			urlPath = urlPath + "?" + raw
		}

		ctx.Set("body", bodyBytes)
		blw := responseLogWriter{responseBody: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = &blw //接口赋值要用地址

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)

		statusCode := ctx.Writer.Status()

		l.Infof("|%d|%v|%s|%s|%s|%s|%s\n",
			statusCode,
			latency,
			clientIP,
			method,
			urlPath,
			strings.ReplaceAll(string(bodyBytes), "\n", ""),
			blw.responseBody.Bytes(),
		)
	}
}
