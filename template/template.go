/**
 * @Time    :2022/10/11 16:47
 * @Author  :Xiaoyu.Zhang
 */

package extemplate

import (
	"errors"
	"fmt"
	"github.com/melf-xyzh/go-excel/commons"
	"github.com/melf-xyzh/go-excel/constant"
	"github.com/melf-xyzh/go-excel/style"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
	"strconv"
	"strings"
)

// GetTemplate
/**
 *  @Description: 获取文件导入模板
 *  @receiver e
 *  @param tableName
 *  @param tableHead
 *  @return f
 *  @return err
 */
func (e *ExcelConfig) GetTemplate(tableName string, tableHead []string) (f *excelize.File, err error) {
	// 新建Excel文件（工作簿）
	f = excelize.NewFile()
	// 设置工作表名
	if e.SheetName == "" {
		e.SheetName = exconst.DefaultSheetName
	} else {
		// 修改工作表名称
		f.SetSheetName(exconst.DefaultSheetName, e.SheetName)
	}
	// 设置默认列宽
	if e.DefaultColWidth == 0 {
		e.DefaultColWidth = exconst.DefaultColWidth
	}
	// 设置默认行高
	if e.DefaultRowHeight == 0 {
		e.DefaultRowHeight = exconst.DefaultRowHeight
	}
	// 获取有效的最后一列的列名
	// 获取列对应的列名
	name, errName := excelize.ColumnNumberToName(len(tableHead))
	if errName != nil {
		err = errors.New(fmt.Sprintf("获取第%d列对应的列名失败：%s", len(tableHead), errName.Error()))
		return
	}
	// 设置格式
	styleStr := exstyle.NewExStyleStr()
	// 获取对应的StyleId
	styleStrId, errStyle := styleStr.GetStyle(f)
	if errStyle != nil {
		err = errors.New("获取StyleId失败：" + errStyle.Error())
		return
	}
	// 设置格式
	err = f.SetColStyle(e.SheetName, "A:"+name, styleStrId)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置单元格格式 %s:%s 失败：%s", "A", name, err.Error()))
		return
	}
	// 为列设置格式
	err = f.SetColWidth(e.SheetName, "A", name, e.DefaultColWidth)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置 %s:%s 列宽度失败：%s", "A", name, err.Error()))
		return
	}
	// 设置工作表默认格式
	err = f.SetSheetFormatPr(
		e.SheetName,
		excelize.DefaultColWidth(e.DefaultColWidth),
		excelize.DefaultRowHeight(e.DefaultRowHeight),
	)
	if err != nil {
		err = errors.New("设置工作表默认格式失败：" + err.Error())
		return
	}
	// 设置列宽度(特殊)
	if len(e.SpecialColWidth) > 0 {
		for colName, width := range e.SpecialColWidth {
			err = f.SetColWidth(e.SheetName, colName, colName, width)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置 %s 列宽度失败：%s", colName, err.Error()))
				return
			}
		}
	}
	// 设置行高度(特殊)
	if len(e.SpecialRowHeight) > 0 {
		for row, height := range e.SpecialRowHeight {
			err = f.SetRowHeight(e.SheetName, row, height)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置 %d 行高度失败：%s", row, err.Error()))
				return
			}
		}
	}
	// 设置标题
	err = f.SetCellValue(e.SheetName, "A1", tableName)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置标题失败：%s", err.Error()))
		return
	}
	// 设置表头
	err = f.SetSheetRow(e.SheetName, "A2", &tableHead)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置表头失败：%s", err.Error()))
		return
	}

	// 设置格式
	if len(e.Style) > 0 {
		for cellStr, style := range e.Style {
			cells := strings.Split(cellStr, ",")
			if len(cells) < 2 {
				err = errors.New(fmt.Sprintf("Style Key %s 格式不合法", cellStr))
				return
			}
			hCell, vCell := cells[0], cells[1]
			// 获取对应的StyleId
			styleId, errStyle := style.GetStyle(f)
			if errStyle != nil {
				err = errors.New("获取StyleId失败：" + errStyle.Error())
				return
			}
			// 设置格式
			err = f.SetCellStyle(e.SheetName, hCell, vCell, styleId)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置单元格格式 %s:%s 失败：%s", hCell, vCell, err.Error()))
				return
			}
		}
	} else {
		style := exstyle.NewExStyle(
			exconst.DefaultFontFamily, exconst.DefaultFontSize,
			exconst.DefaultHorizontalAlign, exconst.DefaultVerticalAlign,
		)
		styleId, errStyle := style.GetStyle(f)
		if errStyle != nil {
			err = errors.New("获取StyleId失败：" + errStyle.Error())
			return
		}
		hCell, vCell := "A1", name+"4"
		// 设置格式
		err = f.SetCellStyle(e.SheetName, hCell, vCell, styleId)
		if err != nil {
			err = errors.New(fmt.Sprintf("设置 %s:%s 单元格格式失败：%s", hCell, vCell, err.Error()))
			return
		}
	}
	// 合并单元格
	if len(e.MergeCell) > 0 {
		for hCell, vCell := range e.MergeCell {
			// 合并单元格
			err = f.MergeCell(e.SheetName, hCell, vCell)
			if err != nil {
				err = errors.New(fmt.Sprintf("合并 %s:%s 失败：%s", hCell, vCell, err.Error()))
				return
			}
		}
	} else {
		hCell, vCell := "A1", name+"1"
		// 合并单元格
		err = f.MergeCell(e.SheetName, hCell, vCell)
		if err != nil {
			err = errors.New(fmt.Sprintf("合并 %s:%s 失败：%s", hCell, vCell, err.Error()))
			return
		}
	}
	e.f = f
	return
}

// GetTemplateByStruct
/**
 *  @Description: 通过结构体获取文件导入模板
 *  @receiver e
 *  @param tableName
 *  @param data
 *  @return f
 *  @return err
 */
func (e *ExcelConfig) GetTemplateByStruct(tableName string, data interface{}) (f *excelize.File, err error) {
	// 解析结构体
	tableHead := make([]string, 0)
	m := make(map[int]ExcelTag)
	tableHead, m, err = parse(data, 1)
	if err != nil {
		return
	}
	if len(m) > 0 {
		for k, v := range m {
			// 获取列对应的列名
			name, errName := excelize.ColumnNumberToName(k)
			if errName != nil {
				err = errors.New(fmt.Sprintf("获取第%d列对应的列名失败：%s", len(tableHead), errName.Error()))
				return
			}
			if e.SpecialColWidth == nil {
				e.SpecialColWidth = map[string]float64{}
			}
			e.SpecialColWidth[name] = v.Width
		}
	}
	// 创建模板
	f, err = e.GetTemplate(tableName, tableHead)
	return
}

// ExportFile
/**
 *  @Description: 导出文件
 *  @receiver e
 *  @param filePath
 *  @return err
 */
func (e *ExcelConfig) ExportFile(filePath string) (err error) {
	if e.f == nil {
		err = errors.New("excelize.File 为 nil")
		return
	}
	// 判断保存路径对应的文件夹是否存在
	ok, _ := excommons.PathExists(filePath)
	if !ok {
		// 创建多次文件夹
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			err = errors.New("创建文件夹失败：" + err.Error())
			return err
		}
	}
	// 根据指定路径保存文件
	if e.FileName == "" {
		e.FileName = exconst.DefaultExcelFileName
	} else {
		// 扩展名
		ext := path.Ext(e.FileName)
		if ext == "" {
			e.FileName += exconst.ExcelExt
		} else if ext != exconst.ExcelExt {
			err = errors.New("错误的文件扩展名：" + ext)
			return
		}
	}
	// 导出文件
	err = e.f.SaveAs(path.Join(filePath, e.FileName))
	return
}

// GetLadderTemplate
/**
 *  @Description: 生成阶梯导入模板
 *  @receiver e
 *  @param tableName
 *  @param firstCase
 *  @param tableHead
 *  @return f
 *  @return err
 */
func (e *ExcelConfig) GetLadderTemplate(firstCase string, tableHead []string) (f *excelize.File, err error) {
	// 新建Excel文件（工作簿）
	f = excelize.NewFile()
	// 设置工作表名
	if e.SheetName == "" {
		e.SheetName = exconst.DefaultSheetName
	} else {
		// 修改工作表名称
		f.SetSheetName(exconst.DefaultSheetName, e.SheetName)
	}
	// 设置默认列宽
	if e.DefaultColWidth == 0 {
		e.DefaultColWidth = exconst.DefaultColWidth
	}
	// 设置默认行高
	if e.DefaultRowHeight == 0 {
		e.DefaultRowHeight = exconst.DefaultRowHeight
	}
	// 获取有效的最后一列的列名
	// 获取列对应的列名
	name, errName := excelize.ColumnNumberToName(len(tableHead) + 1)
	if errName != nil {
		err = errors.New(fmt.Sprintf("获取第%d列对应的列名失败：%s", len(tableHead), errName.Error()))
		return
	}
	// 设置格式
	styleStr := exstyle.NewExStyleStr()
	// 获取对应的StyleId
	styleStrId, errStyle := styleStr.GetStyle(f)
	if errStyle != nil {
		err = errors.New("获取StyleId失败：" + errStyle.Error())
		return
	}
	// 设置格式
	err = f.SetColStyle(e.SheetName, "A:"+name, styleStrId)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置单元格格式 %s:%s 失败：%s", "A"+strconv.Itoa(len(tableHead)+1), name+strconv.Itoa(len(tableHead)+1), err.Error()))
		return
	}
	// 为列设置格式
	err = f.SetColWidth(e.SheetName, "A", name, e.DefaultColWidth)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置 %s:%s 列宽度失败：%s", "A", name, err.Error()))
		return
	}
	// 设置工作表默认格式
	err = f.SetSheetFormatPr(
		e.SheetName,
		excelize.DefaultColWidth(e.DefaultColWidth),
		excelize.DefaultRowHeight(e.DefaultRowHeight),
	)
	if err != nil {
		err = errors.New("设置工作表默认格式失败：" + err.Error())
		return
	}
	// 设置列宽度(特殊)
	if len(e.SpecialColWidth) > 0 {
		for colName, width := range e.SpecialColWidth {
			err = f.SetColWidth(e.SheetName, colName, colName, width)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置 %s 列宽度失败：%s", colName, err.Error()))
				return
			}
		}
	}
	// 设置行高度(特殊)
	if len(e.SpecialRowHeight) > 0 {
		for row, height := range e.SpecialRowHeight {
			err = f.SetRowHeight(e.SheetName, row, height)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置 %d 行高度失败：%s", row, err.Error()))
				return
			}
		}
	}
	// 设置首格
	err = f.SetCellValue(e.SheetName, "A1", firstCase)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置标题失败：%s", err.Error()))
		return
	}
	// 设置首行
	err = f.SetSheetRow(e.SheetName, "B1", &tableHead)
	if err != nil {
		err = errors.New(fmt.Sprintf("设置表头失败：%s", err.Error()))
		return
	}
	// 设置首列
	//err = f.SetSheetCol(e.SheetName, "A2", &tableHead)
	//if err != nil {
	//	err = errors.New(fmt.Sprintf("设置首列失败：%s", err.Error()))
	//	return
	//}
	for i, s := range tableHead {
		err = f.SetCellValue(e.SheetName, "A"+strconv.Itoa(i+2), s)
		if err != nil {
			err = errors.New(fmt.Sprintf("设置首列值失败：%s", err.Error()))
			return
		}
	}
	// 设置格式
	if len(e.Style) > 0 {
		for cellStr, style := range e.Style {
			cells := strings.Split(cellStr, ",")
			if len(cells) < 2 {
				err = errors.New(fmt.Sprintf("Style Key %s 格式不合法", cellStr))
				return
			}
			hCell, vCell := cells[0], cells[1]
			// 获取对应的StyleId
			styleId, errStyle := style.GetStyle(f)
			if errStyle != nil {
				err = errors.New("获取StyleId失败：" + errStyle.Error())
				return
			}
			// 设置格式
			err = f.SetCellStyle(e.SheetName, hCell, vCell, styleId)
			if err != nil {
				err = errors.New(fmt.Sprintf("设置单元格格式 %s:%s 失败：%s", hCell, vCell, err.Error()))
				return
			}
		}
	} else {
		style := exstyle.NewExStyle(
			exconst.DefaultFontFamily, exconst.DefaultFontSize,
			exconst.DefaultHorizontalAlign, exconst.DefaultVerticalAlign,
		)
		styleId, errStyle := style.GetStyle(f)
		if errStyle != nil {
			err = errors.New("获取StyleId失败：" + errStyle.Error())
			return
		}
		hCell, vCell := "A1", name+strconv.Itoa(len(tableHead)+2)
		// 设置格式
		err = f.SetCellStyle(e.SheetName, hCell, vCell, styleId)
		if err != nil {
			err = errors.New(fmt.Sprintf("设置 %s:%s 单元格格式失败：%s", hCell, vCell, err.Error()))
			return
		}
	}
	// 合并单元格
	if len(e.MergeCell) > 0 {
		for hCell, vCell := range e.MergeCell {
			// 合并单元格
			err = f.MergeCell(e.SheetName, hCell, vCell)
			if err != nil {
				err = errors.New(fmt.Sprintf("合并 %s:%s 失败：%s", hCell, vCell, err.Error()))
				return
			}
		}
	}
	e.f = f
	return
}
