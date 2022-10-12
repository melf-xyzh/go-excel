# go-excel
Excel数据读写的简易封装

## 安装

```go
go get github.com/melf-xyzh/go-excel
```

## 生成导入模板

生成一个默认的Excel导入模板

```go
// 创建一个excelConfig（每个Excel文件需要一个）
exTemp := extemplate.ExcelConfig{}
```

使用自定义配置生成Excel导入版本

```go
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
```

生成文件对象

```go
// 获取Excel导入模板
file, err := exTemp.GetTemplate("信息导入模板", []string{"第1列", "第2列", "第3列", "第4列", "第5列", "第6列"})
if err != nil {
	panic(err)
}
```

通过结构体生成导入模板

```go
// 定义结构体
type User struct {
	Id       string `json:"id"         ex:"column:ID;width:10;required;"`
	Name     string `json:"name"       ex:"column:姓名;width:30;required;"`
	NickName string `json:"nickName"   ex:"column:昵称;width:20;required;"`
	Phone    string `json:"phone"      ex:"column:手机号;width:15;required;"`
	Age      int    `json:"age"        ex:"column:年龄;width:10;required;"`
	Sex      string `json:"sex"        ex:"column:性别;width:10;required;select:男、女"`
	School // 允许存在匿名字段
}

type School struct {
	SchoolName    string `ex:"column:学校;width:30;required;"`
	SchoolAddress string `ex:"column:学校地址;width:50;required;"`
}

// 创建一个excelConfig（每个Excel文件需要一个）
exTemp3 := extemplate.ExcelConfig{
	SheetName:        "车辆导入模板",                                          // 工作表名称
	FileName:         "车辆导入模板" + exconst.DefaultExcelFileName + ".xlsx", // 导出后的文件名
}

user := User{
    Id:       "001",
    Name:     "张三",
    NickName: "别人家的孩子",
    Phone:    "123456789",
    Age:      30,
    Sex:      "男",
}

// 获取Excel导入模板
file, err := exTemp3.GetTemplateByStruct("信息导入模板", user)
if err != nil {
    panic(err)
}
```

保存Excel文件

```go
// 保存Excel文件
err = exTemp.ExportFile("fs")
if err != nil {
	panic(err)
}
```

## 读取导入数据

待开发

## 导出数据

待开发
