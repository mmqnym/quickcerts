package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"QuickCertS/data"
	"QuickCertS/model"
	"QuickCertS/utils"

	"github.com/gin-gonic/gin"
)

// Provide the client with a certificate(unique key and signature) for app.
//
// The server currently uses the device information provided by the client.
//
// Check the device info structure in model/device_info.go.
func ApplyCertificate(ctx *gin.Context) {
	applyInfo := model.ApplyCertInfo{}
	err := ctx.ShouldBindJSON(&applyInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid data format.")
		utils.Logger.Error(err.Error())
		return
	}
	
	invalidResponse := map[string]string{"key": "", "signature": ""}

	excludeList := []string{"Note"} // Allow empty fields.

	// Check the data is all not empty except for the fields in the excludeList.
	if !utils.IsValidData(applyInfo, excludeList) {
		ctx.JSON(http.StatusBadRequest, invalidResponse)
		utils.Logger.Error("Invalid data from client.")
		return
	}

	// Check if the SN exists in the database(It's a legal S/N).
	sn_is_exist, err := data.IsSNExist(applyInfo.SerialNumber)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, invalidResponse)
		return
	}

	if !sn_is_exist {
		ctx.JSON(http.StatusBadRequest, invalidResponse)
		utils.Logger.Error(
			fmt.Sprintf("The S/N [%s] does not exist.", applyInfo.SerialNumber),
		)
		return
	}

	// S/N exists, generate a key and a sinature for the device and update it in the database.
	base := fmt.Sprintf("%s&%s&%s&%s&", 
		applyInfo.SerialNumber, applyInfo.BoardProducer, applyInfo.BoardName, applyInfo.MACAddress)
	key, err := utils.GenerateKey(base)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, invalidResponse)
		utils.Logger.Error(err.Error())
		return
	}

	signature, err := utils.SignMessage([]byte(key))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, invalidResponse)
		utils.Logger.Error(err.Error())
		return
	}

	// Update the key corresponding to the SN in the database.
	// If the verification confirms that the key is the same, resend both the key and signature.
	
	if err := data.BindSNWithKey(applyInfo.SerialNumber, key); err != nil {
		if err.Error() == "the s/n does not exist or has already been used" {
			ctx.JSON(http.StatusBadRequest, err.Error())
			utils.Logger.Warn(
				fmt.Sprintf("The S/N [%s] does not exist or has already been used.", applyInfo.SerialNumber),
			)
		} else {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			utils.Logger.Error(err.Error())
		}
		return

	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Sent the certificate to the client.
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"key":       key,
			"signature": signatureBase64,
		},
	)
	utils.Logger.Info(fmt.Sprintf("Successfully updated and sent the key [%s].", key))
}

// Allow users to apply for temporary use permits on devices they have not previously use the app.
func ApplyTemporaryPermit(ctx *gin.Context) {
	applyInfo := model.ApplyTempPermitInfo{}
	err := ctx.ShouldBindJSON(&applyInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid data format.")
		utils.Logger.Error(err.Error())
		return
	}

	invalidResponse := gin.H{
		"status": "expired",
		"remaining_time": 0,
	}

	excludeList := []string{"Note"} // Allow empty fields.

	// Check the data is all not empty except for the fields in the excludeList.
	if !utils.IsValidData(applyInfo, excludeList) {
		ctx.JSON(http.StatusBadRequest, invalidResponse)
		utils.Logger.Error("Invalid data from client.")
		return
	}

	// Generate a key for the device and update it in the database.
	base := fmt.Sprintf("%s&%s&%s&%s&", 
		"_", applyInfo.BoardProducer, applyInfo.BoardName, applyInfo.MACAddress)
	key, err := utils.GenerateKey(base)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, invalidResponse)
		utils.Logger.Error(err.Error())
		return
	}

	remainingTime, err := data.GetTemporaryPermitExpiredTime(key)

	// The given key has not been used yet, or there is an internal server error.
	if err != nil {
		if strings.Contains(err.Error(), "allowed new key") {
			// Add new key to temporary permit table.
			utils.Logger.Info(err.Error()) // Allowed new key: xxx
			remainingTime, err = data.AddTemporaryPermit(key)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, invalidResponse)
				utils.Logger.Error(err.Error())
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"status": "activated",
					"remaining_time": remainingTime,
				})

				utils.Logger.Info(
					fmt.Sprintf("Authorized [%s] temporary use of the product remaining [%d s].", key, remainingTime),
				)
			}
			
		} else {
			ctx.JSON(http.StatusInternalServerError, invalidResponse)
			utils.Logger.Error(err.Error())
		}
		return
	}

	// Return the remaining valid time.
	if remainingTime > 0 {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"status": "activated",
				"remaining_time": remainingTime,
			},
		)
		utils.Logger.Info(
			fmt.Sprintf("Authorized [%s] temporary use of the product remaining [%d s].", key, remainingTime),
		)
	} else {
		ctx.JSON(http.StatusOK, invalidResponse)
		utils.Logger.Info(
			fmt.Sprintf("The authorization for [%s] to use the product has expired.", key),
		)
	}
	
}