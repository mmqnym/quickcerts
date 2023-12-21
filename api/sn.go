package api

import (
	"fmt"
	"net/http"

	"QuickCertS/data"
	"QuickCertS/model"
	"QuickCertS/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Add serial number to the database, only requests with valid tokens are allowed.
//
// @Summary Create serial number to the database
// @Description Create serial number by providing the serial number and the reason. only requests with valid tokens are allowed.
// @Tags SN
// @Accept json
// @Produce json
// @Param X-RunTime-Code header string false "Security code for admin access. Check path_to_qcs/configs/server.toml for more information."
// @Param X-Access-Token header string false "Security token for admin access. This value is set in path_to_qcs/configs/allowlist.toml."
// @Param snInfo body model.SNInfo true "Serial number information"
// @Success 200 {object} model.CreateSNResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /sn/create [post]
func CreateSN(ctx *gin.Context) {
	creationInfo := model.SNInfo{}
	err := ctx.ShouldBindJSON(&creationInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid data format."})
		utils.Record(logrus.ErrorLevel, err.Error())
		return
	}

	if err := data.AddNewSN(creationInfo.SerialNumber); err != nil {
		if err.Error() == "the s/n already exists" {
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "The S/N already exists."})
			utils.Record(logrus.WarnLevel, fmt.Sprintf("The S/N [%s] already exists.", creationInfo.SerialNumber))
		} else {
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			utils.Record(logrus.ErrorLevel, err.Error())
		}
	} else {
		ctx.JSON(
			http.StatusOK,
			model.CreateSNResponse{Msg: "Successfully uploaded a new S/N.", SerialNumber: creationInfo.SerialNumber},
		)

		utils.Record(logrus.InfoLevel, 
			fmt.Sprintf("Successfully uploaded a new S/N [%s] with reason (%s).",
				creationInfo.SerialNumber, creationInfo.Reason),
		)
	}
}

// Generate serial number(s) to the database, only requests with valid tokens are allowed.
//
// @Summary Generate serial number(s) to the database
// @Description Generate serial number(s) by providing the count and the reason. only requests with valid tokens are allowed.
// @Tags SN
// @Accept json
// @Produce json
// @Param X-RunTime-Code header string false "Security code for admin access. Check path_to_qcs/configs/server.toml for more information."
// @Param X-Access-Token header string false "Security token for admin access. This value is set in path_to_qcs/configs/allowlist.toml."
// @Param snInfo body model.SNsInfo true "Serial number(s) information"
// @Success 200 {object} model.CreateSNResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /sn/generate [post]
func GenerateSN(ctx *gin.Context) {
	generateSNInfo := model.SNsInfo{}
	err := ctx.ShouldBindJSON(&generateSNInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid data format."})
		utils.Record(logrus.ErrorLevel, err.Error())
		return
	}

	snList := []string{}

	if generateSNInfo.Count <= 0 {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "The count must be greater than 0."})
		utils.Record(logrus.WarnLevel, fmt.Sprintf("Invalid count(<=0) [%d].", generateSNInfo.Count))
		return
	}

	for i := 0; i < generateSNInfo.Count; i++ {
		sn, _ := utils.GenerateSN()
		snList = append(snList, sn)
	}

	// Insert snList into database.
	err = data.AddNewSNs(snList)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		utils.Record(logrus.ErrorLevel, err.Error())
	} else {
		msg := fmt.Sprintf("Successfully uploaded new S/N (%d) with reason (%s).",
			generateSNInfo.Count, generateSNInfo.Reason)
		utils.Record(logrus.InfoLevel, msg)
		for _, sn := range snList {
			utils.Record(logrus.InfoLevel, fmt.Sprintf("[%s]", sn))
		}

		ctx.JSON(http.StatusOK, model.GenerateSNResponse{Msg: msg, SerialNumbers: snList})
	}
}

// Add a note for a serial number, only requests with valid tokens are allowed.
//
// @Summary Update a note for a serial number
// @Description Update a note for a serial number by providing the serial number and the note.
// @Tags SN
// @Accept json
// @Produce json
// @Param X-RunTime-Code header string false "Security code for admin access. Check path_to_qcs/configs/server.toml for more information."
// @Param X-Access-Token header string false "Security token for admin access. This value is set in path_to_qcs/configs/allowlist.toml."
// @Param certNote body model.CertNote true "Serial number and note"
// @Success 200 {object} model.UpdateCertNoteResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /sn/update [post]
func UpdateCertNote(ctx *gin.Context) {
	updateInfo := model.CertNote{}
	err := ctx.ShouldBindJSON(&updateInfo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid data format."})
		utils.Record(logrus.ErrorLevel, err.Error())
		return
	}

	if err := data.UpdateCertNote(updateInfo.SerialNumber, updateInfo.Note); err != nil {
		if err.Error() == "the s/n does not exist" {
			errMsg := fmt.Sprintf("The S/N [%s] does not exist.", updateInfo.SerialNumber)
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: errMsg})
			utils.Record(logrus.WarnLevel, errMsg)
		} else {
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			utils.Record(logrus.ErrorLevel, err.Error())
		}

	} else {
		ctx.JSON(
			http.StatusOK,
			model.UpdateCertNoteResponse{
				Msg:  "Successfully updated the note of specified S/N.",
				Note: updateInfo.Note,
			},
		)
		utils.Record(
			logrus.InfoLevel, 
			"Successfully updated the note of specified S/N.",
		)
	}
}

// Get cert list from the database.
//
// @Summary Get cert list from the database
// @Description Get cert list from the database.
// @Tags SN
// @Accept json
// @Produce json
// @Param X-RunTime-Code header string false "Security code for admin access. Check path_to_qcs/configs/server.toml for more information."
// @Param X-Access-Token header string false "Security token for admin access. This value is set in path_to_qcs/configs/allowlist.toml."
// @Success 200 {object} model.GetAllRecordsResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /sn/get-all [get]
func GetAllRecords(ctx *gin.Context) {
	certList, err := data.GetAllCerts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		utils.Record(logrus.ErrorLevel, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, model.GetAllRecordsResponse{Data: certList})
}

// Get available S/N from the database.
//
// @Summary Get available S/N from the database
// @Description Get available S/N from the database.
// @Tags SN
// @Accept json
// @Produce json
// @Param X-RunTime-Code header string false "Security code for admin access. Check path_to_qcs/configs/server.toml for more information."
// @Param X-Access-Token header string false "Security token for admin access. This value is set in path_to_qcs/configs/allowlist.toml."
// @Success 200 {object} model.GetAvaliableSNResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /sn/get-available [get]
func GetAvaliableSN(ctx *gin.Context) {
	snList, err := data.GetAvaliableSN()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		utils.Record(logrus.ErrorLevel, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, model.GetAvaliableSNResponse{Data: snList})
}
