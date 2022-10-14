/**
 * @Time    :2022/10/12 15:57
 * @Author  :Xiaoyu.Zhang
 */

package extemplate

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// parse
/**
 *  @Description: 解析数据
 *  @param data
 *  @param index 递归索引
 *  @return tableHead
 *  @return exTagMap
 *  @return err
 */
func parse(data interface{}, index int) (tableHead []string, exTagMap map[int]ExcelTag, err error) {
	// 获取结构体实例的反射类型对象
	typeOf := reflect.TypeOf(data)
	if typeOf.Kind() != reflect.Struct {
		err = errors.New("非结构体")
		return
	}
	exTagMap = make(map[int]ExcelTag)
	//// 获取值
	valueOf := reflect.ValueOf(data)
	// 可以获取所有属性
	// 获取结构体字段个数：t.NumField()
	for i := 0; i < typeOf.NumField(); i++ {
		// 取每个字段
		f := typeOf.Field(i)
		// 获取字段的值信息
		switch valueOf.Field(i).Kind() {
		case reflect.String, reflect.Slice, reflect.Array:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fallthrough
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 获取Tag
			exTag := f.Tag.Get("ex")
			if exTag == "" || exTag == "-" {
				continue
			}
			var excelTag ExcelTag
			exs := strings.Split(exTag, ";")
			for _, ex := range exs {
				if strings.Contains(ex, "required") {
					excelTag.Required = true
				} else if strings.Contains(ex, ":") {
					kv := strings.Split(ex, ":")
					if len(kv) >= 2 {
						k, v := kv[0], kv[1]
						switch k {
						case "column":
							excelTag.Column = v
						case "width":
							excelTag.Width, _ = strconv.ParseFloat(v, 64)
						case "select":
							selects := strings.Split(v, "、")
							if len(selects)>0 {
								selectMap := make(map[string]struct{})
								for _, sel := range selects {
									selectMap[sel] = struct{}{}
								}
								excelTag.Select = selectMap
							}
						}
					}
				}
				if excelTag.Column == "" {
					excelTag.Column = f.Name
				}
			}
			tableHead = append(tableHead, excelTag.Column)
			exTagMap[index] = excelTag
			index += 1
		case reflect.Ptr:
			continue
		default:
			// Interface()：获取字段对应的值
			val := valueOf.Field(i).Interface()
			//fmt.Println(val)
			tableHead0, exTagMap0, errDG := parse(val, index)
			if errDG != nil {
				return
			}
			for k, v := range exTagMap0 {
				exTagMap[k] = v
			}
			tableHead = append(tableHead, tableHead0...)
		}
	}
	return
}
