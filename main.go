/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 9:25
* @Description:
 */
package main

import (
	"fmt"
	"math/rand"
	_ "proxy-pool/model"
	"proxy-pool/router"
	_ "proxy-pool/scheduler"
	"proxy-pool/utils" //加载顺序先于其他包
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	config := utils.Configure()
	if err := router.Run(fmt.Sprintf("%v:%v", config.System.Host, config.System.Port)); err != nil {
		panic("router Run error:" + err.Error())
	}
}
