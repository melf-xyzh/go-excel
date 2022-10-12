/**
 * @Time    :2022/10/11 16:53
 * @Author  :Xiaoyu.Zhang
 */

package extemplate

import (
	"github.com/melf-xyzh/go-excel/style"
	"github.com/xuri/excelize/v2"
)

// ExcelConfig 表格参数
type ExcelConfig struct {
	SheetName        string                   // 工作表名称
	FileName         string                   // 文件名
	DefaultColWidth  float64                  // 默认列宽
	DefaultRowHeight float64                  // 默认行高
	SpecialColWidth  map[string]float64       // 特殊列宽
	SpecialRowHeight map[int]float64          // 特殊行高
	Style            map[string]exstyle.Style // 格式（map[左上,右下]exstyle.Style）
	MergeCell        map[string]string        // 需要合并单元格（map[左上]右下）
	f                *excelize.File           // Excel文件对象
}
