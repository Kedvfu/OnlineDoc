package models

import "time"

type ExcelData struct {
	ExcelCells  *map[int]*map[int]ExcelCell `json:"excelCells"` // map[row][column]ExcelCell
	OnlineUsers *[]CellHistory              `json:"-"`
}

type ExcelCell struct {
	Content string `json:"content"`
	Style   Style  `json:"style"`
}
type ReceivedExcelCell struct {
	Row     int    `json:"row"`
	Column  int    `json:"column"`
	Content string `json:"content"`
	Style   Style  `json:"style"`
}
type Style struct {
}
type CellHistory struct {
	UserId            int       `json:"userId"`
	Time              time.Time `json:"time"`
	ReceivedExcelCell ReceivedExcelCell
}

func GetEmptyExcelData() *ExcelData {
	onlineUsers := make([]CellHistory, 0)
	excelData := ExcelData{}
	emptyExcelCell := ExcelCell{
		Content: "",
		Style:   Style{},
	}
	excelCellRow := make(map[int]ExcelCell)
	excelCellRow[0] = emptyExcelCell
	excelCellRows := make(map[int]*map[int]ExcelCell)
	excelCellRows[0] = &excelCellRow
	excelData.ExcelCells = &excelCellRows
	excelData.OnlineUsers = &onlineUsers

	return &excelData
}
func (excelData *ExcelData) UpdateExcelCell(row int, column int, content string, style Style, userId int) {
	if excelData.OnlineUsers == nil {
		newCellHistory := make([]CellHistory, 0)
		excelData.OnlineUsers = &newCellHistory
	}
	*excelData.OnlineUsers = append(*excelData.OnlineUsers, CellHistory{
		UserId: userId,
		Time:   time.Now(),
		ReceivedExcelCell: ReceivedExcelCell{
			Row:     row,
			Column:  column,
			Content: content,
			Style:   style,
		},
	})

	newExcelCell := ExcelCell{
		Content: content,
		Style:   style,
	}

	if excelData.ExcelCells == nil {
		if content == "" {
			return
		}
		newRow := make(map[int]ExcelCell)
		newRow[column] = newExcelCell
		newRows := make(map[int]*map[int]ExcelCell)
		newRows[row] = &newRow
		excelData.ExcelCells = &newRows
	} else {
		if (*excelData.ExcelCells)[row] == nil {
			if content == "" {
				return
			}
			newRow := make(map[int]ExcelCell)
			newRow[column] = newExcelCell
			(*excelData.ExcelCells)[row] = &newRow
		} else {
			if content == "" {
				delete(*(*excelData.ExcelCells)[row], column)
				if len(*(*excelData.ExcelCells)[row]) == 0 {
					delete(*excelData.ExcelCells, row)
				}
				return
			}
			(*(*excelData.ExcelCells)[row])[column] = newExcelCell

		}

	}
}
