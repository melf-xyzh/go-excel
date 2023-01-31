/**
 * @Time    :2022/10/13 9:03
 * @Author  :Xiaoyu.Zhang
 */

package extemplate

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/extrame/xls"
	"github.com/melf-xyzh/go-excel/commons"
	"github.com/melf-xyzh/go-excel/constant"
	"github.com/melf-xyzh/go-excel/model"
	"github.com/storyicon/goetag"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var once excommons.ErrOnce

// LoadHttpExcel
/**
 *  @Description:
 *  @receiver e
 *  @param mf http文件对象
 *  @param data
 *  @param ignoreRows
 *  @param filePath 文件保存路径（不包含文件名）
 *  @param fileName 文件名（需要保持的文件名）
 *  @param funModule 所属功能模块名
 *  @param userId 操作人
 *  @param otherFunc 其他操作函数
 *  @return fileRcd 导入导出记录
 *  @return rows excel数据内容
 *  @return err 错误
 */
func (e *ExcelConfig) LoadHttpExcel(mf multipart.File, data interface{}, ignoreRows int, filePath, fileName, funModule, userId string, otherFunc func()) (fileRcd exmodel.FileRecord, rows [][]string, err error) {
	// 读取文件内容
	bs, err := ioutil.ReadAll(mf)
	if err != nil {
		err = errors.New("读取文件异常：" + err.Error())
		return
	}
	// 关闭数据流
	defer mf.Close()
	// 获取ETag
	etag, _ := nameToEtag(fileName)
	// 创建并保存文件
	err = excommons.CreateFile(filePath, etag, bs)
	if err != nil {
		return
	}
	// Minio相关
	if otherFunc != nil {
		otherFunc()
	}
	// 从Excel文件中读取数据
	rows, err = LoadExcelByStruct(filePath, etag, data, ignoreRows)
	if err != nil {
		return
	}
	// 记录操作相关信息
	fileRcd.ID = excommons.UUID()
	fileRcd.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	fileRcd.FileName = fileName
	fileRcd.FilePath = path.Join(filePath, etag)
	fileRcd.Operation = exconst.IMPORT
	fileRcd.FileFormat = exconst.EXCEL
	fileRcd.FunModule = funModule
	fileRcd.CreateUserId = userId
	if e.DB != nil {
		// 单例初始化数据表
		err = once.Do(func() error {
			return e.DB.AutoMigrate(&exmodel.FileRecord{})
		})
		if err != nil {
			err = errors.New("初始化数据导入表失败：" + err.Error())
			return
		}
		err = e.DB.Create(&fileRcd).Error
		if err != nil {
			err = errors.New("数据库异常：" + err.Error())
			return
		}
	}
	return
}

// nameToEtag
/**
 *  @Description: 将名称转换为eTag
 *  @receiver f
 *  @param vf
 */
func nameToEtag(fileName string) (newFileName string, err error) {
	// ETag
	newFileName, err = goetag.GetEtagByString(fileName)
	if err == nil {
		// 扩展名
		ext := path.Ext(fileName)
		newFileName = newFileName + ext
	}
	return
}

// LoadHttpLadderExcel
/**
 *  @Description:
 *  @receiver e
 *  @param mf http文件对象
 *  @param tableHead 表头
 *  @param ignoreRows 忽略行数
 *  @param ignoreCols 忽略列数
 *  @param filePath 文件保存路径（不包含文件名）
 *  @param fileName 文件名（需要保持的文件名）
 *  @param funModule 所属功能模块名
 *  @param userId 操作人
 *  @param otherFunc 其他操作函数
 *  @return fileRcd 导入导出记录
 *  @return rows excel数据内容
 *  @return err 错误
 */
func (e *ExcelConfig) LoadHttpLadderExcel(mf multipart.File, tableHead []string, ignoreRows, ignoreCols int, filePath, fileName, funModule, userId string, otherFunc func()) (fileRcd exmodel.FileRecord, rows [][]string, err error) {
	// 读取文件内容
	bs, err := ioutil.ReadAll(mf)
	if err != nil {
		err = errors.New("读取文件异常：" + err.Error())
		return
	}
	// 关闭数据流
	defer mf.Close()
	// 获取ETag
	etag, _ := nameToEtag(fileName)
	// 创建并保存文件
	err = excommons.CreateFile(filePath, etag, bs)
	if err != nil {
		return
	}
	// Minio相关
	if otherFunc != nil {
		otherFunc()
	}
	// 从Excel文件中读取数据
	rows, err = LoadLadderExcel(filePath, etag, len(tableHead)+1, ignoreRows, ignoreCols)
	if err != nil {
		return
	}
	// 记录操作相关信息
	fileRcd.ID = excommons.UUID()
	fileRcd.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	fileRcd.FileName = fileName
	fileRcd.FilePath = path.Join(filePath, etag)
	fileRcd.Operation = exconst.IMPORT
	fileRcd.FileFormat = exconst.EXCEL
	fileRcd.FunModule = funModule
	fileRcd.CreateUserId = userId
	if e.DB != nil {
		// 单例初始化数据表
		err = once.Do(func() error {
			return e.DB.AutoMigrate(&exmodel.FileRecord{})
		})
		if err != nil {
			err = errors.New("初始化数据导入表失败：" + err.Error())
			return
		}
		err = e.DB.Create(&fileRcd).Error
		if err != nil {
			err = errors.New("数据库异常：" + err.Error())
			return
		}
	}
	return
}

// LoadExcelByStruct
/**
 *  @Description: 导入Excel数据，通过结构体
 *  @param filePath 文件保存路径（不包含文件名）
 *  @param filename
 *  @param data 结构体
 *  @param ignoreRows 忽略行数（对前n行不进行校验）
 *  @return rows 读取出的数据
 *  @return err 错误
 */
func LoadExcelByStruct(filePath, filename string, data interface{}, ignoreRows int) (rows [][]string, err error) {
	var tableHead []string
	var exTagMap map[int]ExcelTag
	tableHead, exTagMap, err = parse(data, 1)
	if err != nil {
		return nil, err
	}
	rows, err = LoadExcel(filePath, filename, len(tableHead))
	if err != nil {
		return nil, err
	}
	uniqueMap := make(map[string]map[string]struct{}, 0)
	for i, row := range rows {
		fmt.Println(fmt.Sprintf("第 %d 行,值：%s", i, row))
		if i < ignoreRows {
			continue
		}
		for j, rowI := range row {
			tag, ok := exTagMap[j+1]
			if !ok {
				continue
			}
			// 必填校验
			if tag.Required {
				if rowI == "" {
					err = errors.New(fmt.Sprintf("第 %d 行,%s 必填", i, tag.Column))
					return
				}
			}

			// 非空则进行其他校验
			if rowI != "" {
				// 唯一校验
				if tag.unique {
					_, exists := uniqueMap[tag.Column][rowI]
					if exists {
						err = errors.New(fmt.Sprintf("第 %d 行,%s(%s) 重复", i, tag.Column, rowI))
						return
					} else {
						_, existsI := uniqueMap[tag.Column]
						if !existsI {
							uniqueMap[tag.Column] = map[string]struct{}{}
						}
						_, existsI = uniqueMap[tag.Column][rowI]
						if !existsI {
							uniqueMap[tag.Column][rowI] = struct{}{}
						}
					}
				}
				// 枚举校验
				if tag.Select != nil {
					_, okSelect := tag.Select[rowI]
					if !okSelect && rowI != "" {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容不合法", i, tag.Column))
						return
					}
				} else if tag.MultiSelect != nil {
					vs := strings.Split(rowI, ",")
					if len(vs) > 0 {
						unique := make(map[string]struct{})
						for _, v := range vs {
							// 多选判重
							_, okUnique := unique[v]
							if okUnique {
								err = errors.New(fmt.Sprintf("第 %d 行,%s（%s） 包含重复值", i, tag.Column, v))
								return
							}
							unique[v] = struct{}{}

							_, okSelect := tag.MultiSelect[v]
							if !okSelect {
								err = errors.New(fmt.Sprintf("第 %d 行,%s（%s） 内容不合法", i, tag.Column, v))
								return
							}
						}
					}
				}

				// 长度校验
				if tag.lens != nil {
					if tag.lens[0] == tag.lens[1] {
						if utf8.RuneCountInString(rowI) != tag.lens[0] {
							err = errors.New(fmt.Sprintf("第 %d 行,%s 内容长度应为 %d 位", i, tag.Column, tag.lens[0]))
							return
						}
					} else if utf8.RuneCountInString(rowI) < tag.lens[0] {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容长度应大于 %d 位", i, tag.Column, tag.lens[0]))
						return
					} else if utf8.RuneCountInString(rowI) > tag.lens[1] {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容长度应小于 %d 位", i, tag.Column, tag.lens[1]))
						return
					}
				}
				// 数值校验（上限）
				if tag.lt != nil {
					v, errNum := strconv.Atoi(rowI)
					if errNum != nil {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容转换为数字失败:%s", i, tag.Column, errNum.Error()))
						return
					}
					if v >= *tag.lt {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 应小于 %d", i, tag.Column, *tag.lt))
						return
					}
				} else if tag.lte != nil {
					v, errNum := strconv.Atoi(rowI)
					if errNum != nil {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容转换为数字失败:%s", i, tag.Column, errNum.Error()))
						return
					}
					if v > *tag.lte {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 应小于等于 %d", i, tag.Column, *tag.lte))
						return
					}
				}
				// 数值校验（下限）
				if tag.gt != nil {
					v, errNum := strconv.Atoi(rowI)
					if errNum != nil {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容转换为数字失败:%s", i, tag.Column, errNum.Error()))
						return
					}
					if v <= *tag.gt {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 应大于 %d", i, tag.Column, *tag.gt))
						return
					}
				} else if tag.gte != nil {
					v, errNum := strconv.Atoi(rowI)
					if errNum != nil {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 内容转换为数字失败:%s", i, tag.Column, errNum.Error()))
						return
					}
					if v < *tag.gte {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 应大于等于 %d", i, tag.Column, *tag.gte))
						return
					}
				}
				// 正则校验
				if tag.Re != "" {
					r := regexp2.MustCompile(tag.Re, 0)
					okRe, _ := r.MatchString(rowI)
					//if errRe != nil {
					//	err = errors.New(fmt.Sprintf("第 %d 行,%s 正则校验失败", i, tag.Column))
					//	return
					//}
					// 正则校验
					if !okRe {
						err = errors.New(fmt.Sprintf("第 %d 行,%s 正则校验失败", i, tag.Column))
						return
					}
				}
			}
		}
	}
	return
}

// LoadExcel
/**
 *  @Description: 导入Excel
 *  @param filePath 文件路径（不包含文件名）
 *  @param filename 文件名
 *  @param colCount 需要读取的列数
 *  @return rows 读取后的内容
 *  @return err 错误
 */
func LoadExcel(filePath, filename string, colCount int) (rows [][]string, err error) {
	filePath = path.Join(filePath, filename)
	// 获取文件扩展名
	ext := excommons.GetFileExt(filename)
	if ext == exconst.ExcelExt2003 {
		var open *xls.WorkBook
		open, err = xls.Open(filePath, "utf-8")
		if err != nil {
			err = errors.New("读取Excel文件失败：" + err.Error())
			return
		}
		// 获取第一个工作表
		sheet := open.GetSheet(0)
		// 遍历xls文件
		for i := 0; i < int(sheet.MaxRow)+1; i++ {
			xlsRow := sheet.Row(i)
			//colCount := xlsRow.LastCol()
			rowData := make([]string, colCount, colCount)
			for j := 0; j < colCount; j++ {
				v := xlsRow.Col(j)
				v = strings.ReplaceAll(v, "\r\n", "")
				v = strings.ReplaceAll(v, "\n", "")
				v = strings.TrimSpace(v)
				rowData[j] = v
			}
			rows = append(rows, rowData)
		}
	} else if ext == exconst.ExcelExt {
		var f *excelize.File
		// 读取excel文件
		f, err = excelize.OpenFile(filePath)
		if err != nil {
			err = errors.New("读取Excel文件失败：" + err.Error())
			return
		}
		// 关闭文件流
		defer f.Close()
		// 获取第一个工作表
		sheet := f.GetSheetName(0)
		// 读取Excel值
		rowsGet, _ := f.GetRows(sheet)
		h := len(rowsGet)
		for i := 1; i <= h; i++ {
			rowData := make([]string, colCount, colCount)
			for j := 1; j < colCount+1; j++ {
				colName, _ := excelize.ColumnNumberToName(j)
				v, _ := f.GetCellValue(sheet, colName+strconv.Itoa(i))
				v = strings.ReplaceAll(v, "\r\n", "")
				v = strings.ReplaceAll(v, "\n", "")
				v = strings.TrimSpace(v)
				rowData[j-1] = v
			}
			rows = append(rows, rowData)
		}
	} else {
		err = errors.New("暂不支持的文件格式")
		return
	}
	return
}

// LoadLadderExcel
/**
 *  @Description:
 *  @param filePath 文件路径
 *  @param filename 文件名
 *  @param colCount 总列数
 *  @param ignoreRows 忽略行数
 *  @param ignoreCols 忽略列数
 *  @return rows 数据
 *  @return err 错误
 */
func LoadLadderExcel(filePath, filename string, colCount, ignoreRows, ignoreCols int) (rows [][]string, err error) {
	rows, err = LoadExcel(filePath, filename, colCount)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		err = errors.New("未读取到数据")
		return
	}
	if len(rows[0])-ignoreCols != colCount-1 {
		err = errors.New("数据列数与期望值不符！")
		return
	}
	newRows := make([][]string, colCount-ignoreRows, colCount-ignoreRows)
	for i, row := range rows {
		log.Println(fmt.Sprintf("第 %d 行,值：%s", i, row))
		if i < ignoreRows {
			continue
		}
		newRows[i-ignoreRows] = row[ignoreCols:]
	}
	rows = newRows
	return
}
