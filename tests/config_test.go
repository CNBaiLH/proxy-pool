/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 9:57
* @Description:
 */
package tests

import (
	"proxy-pool/utils"
	"testing"
)

func TestConfigure(t *testing.T){
	c := utils.Configure()
	cc := c.Database
	t.Logf("%+v",cc)
}
