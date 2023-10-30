package main

import (
	"QuickCertS/api"
	cfg "QuickCertS/configs"
	"QuickCertS/data"
	"QuickCertS/middleware"
	"QuickCertS/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.Logger.Info("Starting the server...")
	utils.Logger.Info("Connecting the database...")

	data.ConnectDB()
	defer data.DisconnectDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AccessLogger())
	allowedIPs := cfg.SERVER_CONFIG.ALLOWED_IPs
	clientAuthTokens := cfg.SERVER_CONFIG.CLIENT_AUTH_TOKEN

	// For client.
	router.POST("/api/apply/cert", middleware.ClientAccessAuth(clientAuthTokens...), api.ApplyCertificate)
	router.POST("/api/apply/temp-permit", middleware.ClientAccessAuth(clientAuthTokens...), api.ApplyTemporaryPermit)

	// For admin.
	router.POST("/api/sn/update", middleware.IPAddressAuth(allowedIPs...), api.UpdateSN)
	router.POST("/api/sn/generate", middleware.IPAddressAuth(allowedIPs...), api.GenerateSN)
	// router.POST("/api/sn/delete", IPAddressAuth(allowedIPs...), api.DeleteSN)
	// router.GET("/api/sn/get-available", IPAddressAuth(allowedIPs...), api.GetAvaliableSN)
	// router.GET("/api/sn/get-all", IPAddressAuth(allowedIPs...), api.GetAllKeyRecords)
	
	utils.Logger.Info("Server initialization is complete and the service will be starting soon...")

	if !cfg.SERVER_CONFIG.USE_TLS {
		run(router)

	} else if cfg.SERVER_CONFIG.TLS_CERT_PATH == "" || cfg.SERVER_CONFIG.TLS_KEY_PATH == "" {
		utils.Logger.Fatal("TLS_CERT_PATH or TLS_KEY_PATH is empty. Please fill in the configs file.")

	} else if cfg.SERVER_CONFIG.SERVE_BOTH {
		runBoth(router)
		
	} else {
		runTLS(router)
	}
}

func run(router *gin.Engine) {
	httpServer := &http.Server{
		Addr:    cfg.SERVER_CONFIG.PORT,
		Handler: router,
		IdleTimeout: cfg.SERVER_CONFIG.KEEP_ALIVE_TIMEOUT,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Failed to start the server. Due to: " + err.Error())
		}
	}()

	utils.WaitForShutdown(httpServer)
}

func runTLS(router *gin.Engine) {
	httpsServer := &http.Server{
		Addr:    cfg.SERVER_CONFIG.TLS_PORT,
		Handler: router,
		IdleTimeout: cfg.SERVER_CONFIG.KEEP_ALIVE_TIMEOUT,
	}

	httpsServer.SetKeepAlivesEnabled(false)

	go func() {
		if err := httpsServer.ListenAndServeTLS(
			cfg.SERVER_CONFIG.TLS_CERT_PATH, 
			cfg.SERVER_CONFIG.TLS_KEY_PATH,
			);
			err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Failed to start the server. Due to: " + err.Error())
		}
	}()

	utils.WaitForShutdown(httpsServer)
}

func runBoth(router *gin.Engine) {
	httpServer := &http.Server{
		Addr:    cfg.SERVER_CONFIG.PORT,
		Handler: router,
		IdleTimeout: cfg.SERVER_CONFIG.KEEP_ALIVE_TIMEOUT,
	}

	httpsServer := &http.Server{
		Addr:    cfg.SERVER_CONFIG.TLS_PORT,
		Handler: router,
		IdleTimeout: cfg.SERVER_CONFIG.KEEP_ALIVE_TIMEOUT,
	}
	
	httpServer.SetKeepAlivesEnabled(false)
	httpsServer.SetKeepAlivesEnabled(false)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Failed to start the server. Due to: " + err.Error())
		}
	}()

	go func() {
		if err := httpsServer.ListenAndServeTLS(
			cfg.SERVER_CONFIG.TLS_CERT_PATH, 
			cfg.SERVER_CONFIG.TLS_KEY_PATH,
			);
			err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Failed to start the server. Due to: " + err.Error())
		}
	}()

	utils.WaitForShutdown(httpServer)
	utils.WaitForShutdown(httpsServer)
}