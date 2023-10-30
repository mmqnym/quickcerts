// Description: Custom Middleware for QCS's Gin Framework

package middleware

import (
	"net"
	"net/http"
	"time"

	"QuickCertS/utils"

	"github.com/gin-gonic/gin"
)

// Override the default logger of Gin Framework.
func AccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		latency := time.Since(start)

		qcsOctx := &utils.QCSExtractGINCtx{
			StatusCode: ctx.Writer.Status(), 
			Latency: latency, 
			ClientIP: ctx.ClientIP(), 
			Method: ctx.Request.Method, 
			FullPath: ctx.FullPath(),
		}

		utils.OverwriteGinLog(qcsOctx)
	}
}

// Middleware for IP authentication.
func IPAddressAuth(allowedIPs ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid format."})
			return
		}

		for _, allowedIP := range allowedIPs {
			if clientIP == allowedIP || allowedIP == "*" {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Request."})
	}
}

func ClientAccessAuth(clientAuthToken ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqToken := ctx.GetHeader("X-Access-Token")
		// var applyInfo model.ApplyInfo
		// err := ctx.Bind(&applyInfo)

		if reqToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Request."})
			return
		}

		for _, allowedToken := range clientAuthToken {
			if reqToken == allowedToken || allowedToken == "" {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Request."})
	}
}