/**
 * @Time    :2022/10/11 16:41
 * @Author  :Xiaoyu.Zhang
 */

package main

import (
	"github.com/melf-xyzh/go-excel/constant"
	"github.com/melf-xyzh/go-excel/style"
	"github.com/melf-xyzh/go-excel/template"
	"log"
)

func main() {
	// 创建一个excelConfig（每个Excel文件需要一个）
	exTemp := extemplate.ExcelConfig{}
	// 获取Excel导入模板
	file, err := exTemp.GetTemplate("信息导入模板", []string{"第1列", "第2列", "第3列", "第4列", "第5列", "第6列"})
	if err != nil {
		panic(err)
	}

	err = exTemp.ExportFile("fs")
	if err != nil {
		panic(err)
	}
	log.Println("保存成功")

	exTemp2 := extemplate.ExcelConfig{
		SheetName:        "车辆导入模板",      // 工作表名称
		FileName:         "车辆导入模板.xlsx", // 导出后的文件名
		DefaultColWidth:  40,            // 默认列宽
		DefaultRowHeight: 300,           // 默认行高（无效）
		SpecialColWidth: map[string]float64{ // 特殊列宽
			"B": 10,
		},
		SpecialRowHeight: map[int]float64{ // 特殊行高
			1: 50,
		},
		Style: map[string]exstyle.Style{ // 自定义单元格格式
			"A1,F3": exstyle.NewExStyle(exconst.DefaultFontFamily, exconst.DefaultFontSize, exconst.DefaultHorizontalAlign, exconst.DefaultVerticalAlign),
			"A1,F1": exstyle.NewExStyle("黑体", 50, "center", "center"),
		},
		MergeCell: map[string]string{ // 需要合并的单元格
			"A1": "D3",
		},
	}
	_, err = exTemp2.GetTemplate("车辆导入模板", []string{"第1列", "第2列", "第3列", "第4列", "第5列", "第6列"})
	if err != nil {
		panic(err)
	}
	err = exTemp2.ExportFile("fs")
	if err != nil {
		panic(err)
	}
	log.Println("保存成功")
}
