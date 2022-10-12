/**
 * @Time    :2022/10/11 17:19
 * @Author  :Xiaoyu.Zhang
 */

package exconst

import "time"

const (
	// ExcelExt 文件扩展名（2007之后版本）
	ExcelExt = ".xlsx"
	// ExcelExt2003 文件扩展名(2007之前版本)
	ExcelExt2003 = ".xls"
)

const (
	// DefaultSheetName 默认工作表名
	DefaultSheetName = "Sheet1"
	// DefaultColWidth 默认列宽
	DefaultColWidth = 25.00
	// DefaultRowHeight 默认行高
	DefaultRowHeight = 20.00
	// DefaultFontFamily 默认字体
	DefaultFontFamily = "宋体"
	// DefaultFontSize 默认字号
	DefaultFontSize = 20.00
	// DefaultHorizontalAlign 默认水平对齐方式
	DefaultHorizontalAlign = "center"
	// DefaultVerticalAlign 默认垂直对齐方式
	DefaultVerticalAlign = "center"
)

var(
	DefaultExcelFileName = time.Now().Format("20060102150405.xlsx")
)