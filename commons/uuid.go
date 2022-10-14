/**
 * @Time    :2022/10/13 11:31
 * @Author  :Xiaoyu.Zhang
 */

package excommons

import (
	"github.com/google/uuid"
	"strings"
)

// UUID 生成uuid
func UUID() string {
	id :=uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}
