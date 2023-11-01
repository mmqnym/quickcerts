package api

import (
	"fmt"
	"net/http"

	"QuickCertS/data"
	"QuickCertS/model"
	"QuickCertS/utils"

	"github.com/gin-gonic/gin"
)

// Add serial number to the database, only requests with valid tokens are allowed.
func UpdateSN(ctx *gin.Context) {
	updateInfo := model.SNInfo{}
	err := ctx.ShouldBindJSON(&updateInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid data format.")
		utils.Logger.Error(err.Error())
		return
	}
	
	if err := data.AddNewSN(updateInfo.SerialNumber); err != nil {
		if err.Error() == "the S/N already exists" {
			ctx.JSON(http.StatusBadRequest, err.Error())
			utils.Logger.Warn(
				fmt.Sprintf("The S/N [%s] already exists, From [%s]", updateInfo.SerialNumber, ctx.RemoteIP()),
			)
		} else {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			utils.Logger.Error(err.Error())
		}
	} else {
		ctx.JSON(http.StatusOK, "Successfully uploaded a new S/N.")
		utils.Logger.Info(
			fmt.Sprintf("Successfully uploaded a new S/N [%s], From [%s]", updateInfo.SerialNumber, ctx.RemoteIP()),
		)
	}
}

// Generate serial number(s) to the database, only requests with valid tokens are allowed.
func GenerateSN(ctx *gin.Context) {
	generateSNInfo := model.SNsInfo{}
	err := ctx.ShouldBindJSON(&generateSNInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid data format.")
		utils.Logger.Error(err.Error())
		return
	}

	snList := []string{}

	if generateSNInfo.Count <= 0 {
		ctx.JSON(http.StatusBadRequest, "The count must be greater than 0.")
		utils.Logger.Warn(fmt.Sprintf("Invalid count(<=0) [%d], From [%s]", generateSNInfo.Count, ctx.RemoteIP()))
		return
	}

	for i := 0; i < generateSNInfo.Count; i++ {
		sn, _ := utils.GenerateSN()
		snList = append(snList, sn)
	}

	// Insert snList into database.
	err = data.AddNewSNs(snList)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		utils.Logger.Error(err.Error())
	} else {
		msg := fmt.Sprintf("Successfully uploaded new S/N (%d), From [%s]", generateSNInfo.Count, ctx.RemoteIP())
		utils.Logger.Info(msg)
		for _, sn := range snList {
			utils.Logger.Info(fmt.Sprintf("[%s]", sn))
		}

		ctx.JSON(http.StatusOK, msg)
	}
}

// Get cert list from the database.
func GetAllSN(ctx *gin.Context) {
	snList, err := data.GetAllCerts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		utils.Logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, snList)
}

// Get available S/N from the database.
func GetAvaliableSN(ctx *gin.Context) {
	snList, err := data.GetAvaliableSN()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		utils.Logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, snList)
}
