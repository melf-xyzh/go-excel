/**
 * @Time    :2022/10/13 9:32
 * @Author  :Xiaoyu.Zhang
 */

package excommons

import (
	"errors"
	"os"
	"path"
)

// PathExists
/**
 *  @Description: 判断路径是否存在
 *  @param path
 *  @return bool
 *  @return error
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	// IsNotExist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// CreateFile
/**
 *  @Description: 创建危机
 *  @param filePath
 *  @param fileName
 *  @param fileContent
 *  @return err
 */
func CreateFile(filePath, fileName string, fileContent []byte) (err error) {
	// 判断保存路径对应的文件夹是否存在
	ok, _ := PathExists(filePath)
	if !ok {
		// 创建多次文件夹
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			err = errors.New("创建文件夹失败：" + err.Error())
			return
		}
	}
	// 拼接全路径
	filePath = path.Join(filePath, fileName)
	// 创建临时文件
	tmp, err := os.Create(filePath)
	if err != nil {
		err = errors.New("创建临时文件异常：" + err.Error())
		return
	}
	// 关闭文件流
	defer tmp.Close()
	// 写文件
	_, err = tmp.Write(fileContent)
	if err != nil {
		err = errors.New("写文件异常：" + err.Error())
		return
	}
	return
}

// GetFileExt
/**
 *  @Description: 获取文件扩展名
 *  @param fileName
 *  @return ext
 */
func GetFileExt(fileName string) (ext string) {
	// 获取文件后缀
	ext = path.Ext(fileName)
	return
}