/**
* @Author: Lanhai Bai
* @Date: 2021/8/26 10:15
* @Description:
 */
package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"proxy-pool/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

//1秒钟才能访问一次
type _IPKuaiDaili struct {
	progress int
	SiteName string
	maxPage  int64
}

var kuaidaili *_IPKuaiDaili
var onceKuaidaili sync.Once

func NewIPKuaiDailiInstance(opts ...func(o *spiderOption)) *_IPKuaiDaili {
	if kuaidaili == nil {
		onceKuaidaili.Do(func() {
			opt := &spiderOption{maxPage:6}
			for _,fn := range opts {
				fn(opt)
			}
			kuaidaili = &_IPKuaiDaili{SiteName: "快代理", progress: 0, maxPage: opt.maxPage}
		})
	}
	return kuaidaili
}

func (i *_IPKuaiDaili) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	urls := []requestURL{
		{URL: "https://www.kuaidaili.com/free/inha", Type: "高匿"},
		{URL: "https://www.kuaidaili.com/free/intr", Type: "普通"},
	}
	var crawldata func(u string, baseURL string)
	crawldata = func(u string, baseURL string) {
		utils.Logger().Infof("current url :[%v]", u)
		var dom *goquery.Document
		if dom = getDOMByURL(u); dom == nil {
			utils.Logger().Debugf("IPKuaiDaili CrawlData getDOMByURL url:[%v] dom is empty", u)
			return
		}
		lis := dom.Find("#listnav").Find("li")
		length := lis.Length()
		if length < 2 {
			utils.Logger().Debugf("IPKuaiDaili pagenation data length is: %v", length)
			return
		}
		pagenationString := lis.Eq(length - 2).Children().Text()
		page, _ := strconv.Atoi(pagenationString)
		utils.Logger().Infof("IPKuaiDaili total page: [%v]", page)

		i.Paraser(dom, nil)
		if page > int(i.maxPage) {
			page = int(i.maxPage)
		}

		for j := 2; j <= page; j++ {
			<-time.After(time.Second)
			if d := getDOMByURL(fmt.Sprintf("%v/%v", baseURL, j)); d != nil {
				go i.Paraser(d, nil)
				return
			}
		}
	}

	for _, v := range urls {
		crawldata(v.URL, v.URL)
		<-time.After(time.Second)
	}
}

func (ii *_IPKuaiDaili) Paraser(dom *goquery.Document, params map[string]interface{}) {
	tbody := dom.Find("#list").Find("tbody")
	if tbody == nil {
		return
	}
	trs := tbody.Find("tr")
	length := trs.Length()

	for i := 0; i < length; i++ {
		tds := trs.Eq(i).Find("td")
		ip := replaceSpaceAndOtherMark(tds.Eq(0).Text())
		port := replaceSpaceAndOtherMark(tds.Eq(1).Text())
		isHTTPs := 0
		if strings.Contains(tds.Eq(3).Text(), "HTTPS") {
			isHTTPs = 1
		}
		local := replaceSpaceAndOtherMark(tds.Eq(4).Text())
		ProxiesChannel <- constructProxy(ip,port,isHTTPs,0,local,"未知",ii.SiteName)
	}
}

func (i *_IPKuaiDaili) Progress() int {
	return i.progress
}
func(i *_IPKuaiDaili) ProgressName() string{
	return i.SiteName
}
