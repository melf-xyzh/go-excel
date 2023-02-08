/**
 * @Time    :2022/10/11 17:41
 * @Author  :Xiaoyu.Zhang
 */

package exstyle

import "github.com/xuri/excelize/v2"

type Style struct {
	excelize.Style
}

// NewExStyle
/**
 *  @Description: 创建Excel单元格样式
 *  @param family 字体
 *  @param size 字号
 *  @param horizontal 水平对齐方式（left-向左（缩进）、center-居中、right-靠右（缩进）、fill-填充、justify-两端对齐、centerContinuous-跨列居中、distributed-分散对齐（缩进））
 *  @param vertical 垂直居中方式（top-顶端对齐、center-居中、justify-两端对齐、distributed-分散对齐）
 */
func NewExStyle(family string, size float64, horizontal, vertical string) (excelStyle Style) {
	excelStyle.Font = &excelize.Font{
		Family: family,
		Size:   size,
	}
	excelStyle.Alignment = &excelize.Alignment{
		Horizontal: horizontal,
		Vertical:   vertical,
		WrapText:   true,
	}
	// 设置文件格式为纯文本
	excelStyle.NumFmt = 49
	return
}

// NewExStyleStr
/**
 *  @Description: 设置文件格式为纯文本
 *  @return excelStyle
 */
func NewExStyleStr() (excelStyle Style) {
	excelStyle.NumFmt = 49
	return
}

// GetStyle
/**
 *  @Description: 获取Style对象
 *  @receiver s
 *  @param f
 *  @return style
 *  @return err
 */
func (s *Style) GetStyle(f *excelize.File) (style int, err error) {
	style, err = f.NewStyle(&s.Style)
	return
}

// SetFamily
/**
 *  @Description: 设置字体
 *  @receiver s
 *  @param family
 */
func (s *Style) SetFamily(family string) {
	s.Font.Family = family
}
