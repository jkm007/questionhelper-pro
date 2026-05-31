package excel

import (
	"fmt"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

// ReadExcel 读取 Excel 文件
func ReadExcel(file *multipart.FileHeader) ([][]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("读取 Excel 失败: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("获取行数据失败: %w", err)
	}

	return rows, nil
}

// WriteExcel 写入 Excel 文件
func WriteExcel(headers []string, data [][]string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	index, _ := f.NewSheet(sheetName)

	// 写入表头
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for rowIdx, row := range data {
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	f.SetActiveSheet(index)
	return f, nil
}

// ExportExcel 导出 Excel 文件
func ExportExcel(filename string, headers []string, data [][]string) error {
	f, err := WriteExcel(headers, data)
	if err != nil {
		return err
	}
	return f.SaveAs(filename)
}
