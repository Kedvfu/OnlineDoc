package handlers

import (
	"OnlineDoc/api/sessions"
	"OnlineDoc/models"
	"github.com/gin-gonic/gin"
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
		return
	}

	var receivedExcelCell models.ReceivedExcelCell
	err = context.ShouldBindJSON(&receivedExcelCell)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document content",
		})
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
		return
	}
	refreshExcelData := *(*sessions.ExcelSessions)[documentIdNum].OnlineUsers

	var jsonData map[string]float64
	err = context.ShouldBindJSON(&jsonData)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document content",
		})

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
