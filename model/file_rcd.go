/**
 * @Time    :2022/10/13 9:20
 * @Author  :Xiaoyu.Zhang
 */

package exmodel

// FileRecord 文件导入导出记录表
type FileRecord struct {
	ID           string `json:"id,omitempty"               gorm:"column:id;primary_key;type:varchar(36)"`
	CreateTime   string `json:"createTime,omitempty"       gorm:"column:create_time;index;type:varchar(20)"`
	FileName     string `json:"fileName"                   gorm:"column:file_name;comment:文件名;type:varchar(100);"`
	FilePath     string `json:"filePath"                   gorm:"column:file_path;comment:文件路径;type:text;"`
	Operation    string `json:"operation"                  gorm:"column:operation;comment:文件操作;type:varchar(20);"`
	FileFormat   string `json:"fileFormat"                 gorm:"column:file_format;comment:文件格式;type:varchar(20);"`
	FunModule    string `json:"funModule"                  gorm:"column:fun_module;comment:所属功能模块;type:varchar(255);"`
	CreateUserId string `json:"-"                          gorm:"column:create_user_id;comment:创建人;type:varchar(36);"`
}

// TableName 自定义表名
func (FileRecord) TableName() string {
	return "file_rcd"
}
