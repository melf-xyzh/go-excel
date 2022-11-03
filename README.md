# go-excel
Excel数据读写的简易封装

## 安装

```go
go get -u github.com/melf-xyzh/go-excel
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

生成阶梯模板

```go
// 获取Excel导入模板
_, err := exTemp.GetLadderTemplate("乘法表", []string{"1", "2", "3", "4", "5","6", "7", "8", "9"})
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

通过结构体读入数据

```go
user := User{}
// 通过结构体读取数据，并进行校验
// rows, err := extemplate.LoadExcelByStruct("./", "book1.xlsx", user, 2)
rows, err := extemplate.LoadExcelByStruct("./", "book1.xls", user, 2)
if err != nil {
	panic(err)
}
for _, row := range rows {
	fmt.Println(row)
}
```

常规方式读取数据

```go
// rows, err = extemplate.LoadExcel("./", "book1.xls", 8)
rows, err = extemplate.LoadExcel("./", "book1.xlsx", 8)
if err != nil {
	panic(err)
}
```

读取阶梯模板

```go
rows, err := extemplate.LoadLadderExcel("./","乘法表.xlsx",10,1,1)
if err != nil {
	panic(err)
}
for _, row := range rows {
	fmt.Println(row)
}
```

## 导出数据

待开发

## exTag

| 标签名   | 说明                                                         |
| :------- | :----------------------------------------------------------- |
| column   | 指定导出`Excel`时的列名<br />ex:"column:姓名;"               |
| width    | 指定导出`Excel`时的列宽<br />ex:"width:30;"                  |
| required | 必填<br />ex:"required;"                                     |
| len      | 导入`Excel`时的数据长度校验<br />ex:"len:2-5"（2-5位）<br />ex:"len:3"（固定为3位） |
| re       | 正则校验<br /> ex:"re:^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[A-Z]{1}[A-Z0-9]{4,5}[A-Z0-9挂学警港澳]{1}$" |
| >        | 数值下限校验<br /> ex:">:0"                                  |
| >=       | 数值下限校验<br /> ex:">=:0"                                 |
| <        | 数值上限校验<br /> ex:"<:100"                                |
| <=       | 数值上限校验<br /> ex:"<=:100"                               |

