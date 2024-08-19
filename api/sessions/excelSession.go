package sessions

import "OnlineDoc/models"

var ExcelSessions *map[int]models.ExcelData //map[DocumentID]models.ExcelData

func InitialExcelSessions() {
	excelSessions := make(map[int]models.ExcelData)
	ExcelSessions = &excelSessions
}
