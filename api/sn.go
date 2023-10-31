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

	owner := utils.GetValidTokenOwner(updateInfo.Token)

	if owner == "" {
		ctx.JSON(http.StatusUnauthorized, "")
		utils.Logger.Warn(fmt.Sprintf("Illegal access detected, From [%s]", ctx.RemoteIP()))
		return
	} else {
		utils.Logger.Info(fmt.Sprintf("Admin [%s] login, From [%s]", owner, ctx.RemoteIP()))
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
			fmt.Sprintf("Successfully uploaded a new S/N [%s] by [%s], From [%s]", 
				updateInfo.SerialNumber, owner, ctx.RemoteIP()),
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

	owner := utils.GetValidTokenOwner(generateSNInfo.Token)

	if owner == "" {
		ctx.JSON(http.StatusUnauthorized, "")
		utils.Logger.Warn(fmt.Sprintf("Illegal access detected, From [%s]", ctx.RemoteIP()))
		return
	} else {
		utils.Logger.Info(fmt.Sprintf("Admin [%s] login, From [%s]", owner, ctx.RemoteIP()))
	}

	// Generate sn * count.
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
		msg := fmt.Sprintf("Successfully uploaded new S/N. (%d).", generateSNInfo.Count)
		ctx.JSON(http.StatusOK, msg)
		utils.Logger.Info(fmt.Sprintf("%s by [%s], From [%s]", msg, owner, ctx.RemoteIP()))
		for _, sn := range snList {
			utils.Logger.Info(fmt.Sprintf("[%s]", sn))
		}
	}
}