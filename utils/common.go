/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 9:44
* @Description:
 */
package utils

import (
	"os"
	"path/filepath"
)

func GetWebDir() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("load config error:"+err.Error())
		return
	}
	return // strings.Replace(dir, "\\", "/", -1)
}