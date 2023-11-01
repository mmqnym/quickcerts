// Description: Custom Middleware for QCS's Gin Framework

package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	cfg "QuickCertS/configs"
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
func IPAddressAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid format."})
			return
		}

		for _, allowedIP := range cfg.SERVER_CONFIG.ALLOWED_IPs {
			if clientIP == allowedIP || allowedIP == "*" {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Request."})
	}
}

// Middleware for client authentication.
func ClientAccessAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqToken := ctx.GetHeader("X-Access-Token")

		if reqToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Request."})
			return
		}

		for _, allowedToken := range cfg.SERVER_CONFIG.CLIENT_AUTH_TOKEN {
			if reqToken == allowedToken || allowedToken == "" {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Request."})
	}
}


func AdminAccessAuth(runTimeCode string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqRunTimeCode := ctx.GetHeader("X-Runtime-Code")
		reqToken := ctx.GetHeader("X-Access-Token")

		if cfg.SERVER_CONFIG.USE_RUNTIME_CODE {
			if reqRunTimeCode == "" || reqRunTimeCode != runTimeCode {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Request."})
				return
			}
		}

		if reqToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Request."})
			return
		}

		for _, permission := range cfg.ALLOWEDLIST.PERMISSIONS {
			if reqToken == permission.TOKEN || permission.TOKEN == "" {
				utils.Logger.Info(fmt.Sprintf("Admin [%s] login, From [%s]", permission.NAME, ctx.RemoteIP()))
				ctx.Next()
				return
			}
		}

		utils.Logger.Warn(fmt.Sprintf("Illegal access detected, From [%s]", ctx.RemoteIP()))
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Request."})
	}
}