package handlers

import (
	"OnlineDoc/api/sessions"
	"OnlineDoc/files"
	"OnlineDoc/models"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

func UpdateExcel(context *gin.Context) {

	userId, _ := context.Get("userId")
	userIdNum := userId.(int)

	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)

	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		context.Abort()
		return
	}

	var receivedExcelCell models.ReceivedExcelCell
	err = context.ShouldBindJSON(&receivedExcelCell)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document content",
		})
		context.Abort()
		return
	}
	updateRow := receivedExcelCell.Row
	updateCol := receivedExcelCell.Column
	updateContent := receivedExcelCell.Content
	//updateStyle := jsonData.Style

	updateExcelData := (*sessions.ExcelSessions)[documentIdNum]

	updateExcelData.UpdateExcelCell(updateRow, updateCol, updateContent, models.Style{}, userIdNum)
	context.JSON(200, gin.H{
		"status": 1,
	})
}
func RefreshExcel(context *gin.Context) {

	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		context.Abort()
		return
	}
	refreshExcelData := *(*sessions.ExcelSessions)[documentIdNum].OnlineUsers

	var jsonData map[string]float64
	err = context.ShouldBindJSON(&jsonData)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document content",
		})
		context.Abort()
		return
	}
	userLastRefreshTime := time.UnixMilli(int64(jsonData["timestamp"]))
	cellsNeedToUpdate := make([]models.CellHistory, 0)
	for _, cellHistory := range refreshExcelData {
		if cellHistory.Time.After(userLastRefreshTime) {
			cellsNeedToUpdate = append(cellsNeedToUpdate, cellHistory)
		}
	}

	context.JSON(200, gin.H{
		"cells": cellsNeedToUpdate,
	})
}

func DownloadExcel(context *gin.Context) {

	documentId, _ := context.Get("documentId")
	documentIdNum := documentId.(int)
	var fileBytes []byte

	BlobFile, err := RedisClient.Get("document_download_" + strconv.Itoa(documentIdNum)).Bytes()
	if err != nil {
		file := excelize.NewFile()
		sheetName := "sheet1"
		index, err := file.NewSheet(sheetName)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "unable to create new file",
			})
			context.Abort()
			return
		}
		file.SetActiveSheet(index)
		err = files.WriteExcelCellsToFile(file, sheetName, documentIdNum)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "unable to write to file",
			})
			context.Abort()
			return

		}
		var buf bytes.Buffer
		if err := file.Write(&buf); err != nil {
			context.JSON(500, gin.H{
				"message": "unable to write file to buffer",
			})
			context.Abort()
			return
		}

		context.Header("Content-Length", strconv.Itoa(buf.Len()))
		RedisClient.Set("document_download_"+strconv.Itoa(documentIdNum), buf.Bytes(), 3600*time.Second)

	} else {
		fileBytes = BlobFile
		context.Header("Content-Length", strconv.Itoa(len(fileBytes)))
	}
	fileName := "Excel_" + strconv.Itoa(documentIdNum) + "_" + time.Now().Format("20060102-150405") + ".xlsx"
	context.Header("Content-Disposition", "attachment; filename="+fileName)
	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 将字节缓冲区的内容发送到客户端
	context.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)

}
