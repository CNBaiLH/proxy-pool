/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 11:28
* @Description:
 */
package tests

import (
	"proxy-pool/service"
	"testing"
)

func TestIPYQie(t *testing.T) {
	i := service.NewIPXilaInstance()
	i.CrawlData()
}