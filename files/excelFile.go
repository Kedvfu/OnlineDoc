package files

import (
	"OnlineDoc/api/sessions"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func WriteExcelCellsToFile(file *excelize.File, sheetName string, documentId int) error {
	excelCells := (*sessions.ExcelSessions)[documentId].ExcelCells

	for row, excelRow := range *excelCells {
		for column, excelCell := range *excelRow {
			if row != 0 && column != 0 {
				err := file.SetCellValue(sheetName, GetPositionString(row, column), excelCell.Content)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func GetPositionString(row int, column int) string {
	var columnString string

	for column > 0 {
		column--
		remainder := column % 26
		columnString = string(rune('A'+remainder)) + columnString
		column /= 26
	}
	return columnString + strconv.Itoa(row)
}
