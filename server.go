package main

import (
	"QuickCertS/api"
	cfg "QuickCertS/configs"
	"QuickCertS/data"
	"QuickCertS/middleware"
	"QuickCertS/utils"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var (
	runtimeCode string
	router *gin.Engine
)

func init() {
	utils.Logger.Info("Initializing the server...")

	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AccessLogger())

	if cfg.SERVER_CONFIG.USE_RUNTIME_CODE {
		var err error
		runtimeCode, err = utils.GenerateRunTimeCode()

		if err != nil {
			utils.Logger.Fatal("Failed to generate the run time code. Due to: " + err.Error())
		}

		runtimeCodeMsg := color.HiCyanString("[USE_RUNTIME_CODE] is enabled, Runtime code: ")
		runtimeCodeMsg += color.HiMagentaString("%s", runtimeCode)
		utils.Logger.Info(runtimeCodeMsg)
	}
}

func main() {
	data.ConnectDB()
	defer data.DisconnectDB()
	
	registRoutes()

	if !cfg.SERVER_CONFIG.USE_TLS {
		run(router)

	} else if cfg.SERVER_CONFIG.TLS_CERT_PATH == "" || cfg.SERVER_CONFIG.TLS_KEY_PATH == "" {
		utils.Logger.Fatal("TLS_CERT_PATH or TLS_KEY_PATH is empty. Please fill in the configs file.")
		
	} else {
		runTLS(router)
	}
}

func registRoutes() {
	// For client.
	router.POST("/api/apply/cert", middleware.ClientAccessAuth(), api.ApplyCertificate)
	router.POST("/api/apply/temp-permit", middleware.ClientAccessAuth(), api.ApplyTemporaryPermit)

	// For admin.
	router.POST("/api/sn/update", 
		middleware.IPAddressAuth(), 
		middleware.AdminAccessAuth(runtimeCode), 
		api.UpdateSN,
	)
	router.POST("/api/sn/generate", 
		middleware.IPAddressAuth(), 
		middleware.AdminAccessAuth(runtimeCode), 
		api.GenerateSN,
	)
	// router.POST("/api/sn/delete", IPAddressAuth(), api.DeleteSN)
	router.GET("/api/sn/get-available", 
		middleware.IPAddressAuth(), 
		middleware.AdminAccessAuth(runtimeCode), 
		api.GetAvaliableSN,
	)
	router.GET("/api/sn/get-all", 
		middleware.IPAddressAuth(), 
		middleware.AdminAccessAuth(runtimeCode), 
		api.GetAllSN,
	)
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

	runningMsg := fmt.Sprintf("Server is running in %s mode. listening on port: %s", 
		color.HiCyanString("http"), color.HiCyanString("%s", cfg.SERVER_CONFIG.PORT[1:]))
	utils.Logger.Info(runningMsg)

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

	runningMsg := fmt.Sprintf("Server is running in %s mode. listening on port: %s", 
		color.HiMagentaString("https"), color.HiMagentaString("%s", cfg.SERVER_CONFIG.TLS_PORT[1:]))
	utils.Logger.Info(runningMsg)

	utils.WaitForShutdown(httpsServer)
}