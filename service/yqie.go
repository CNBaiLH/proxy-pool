/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 14:34
* @Description:
 */
package service

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"proxy-pool/utils"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//http://ip.yqie.com/ipproxy.htm
//https http://ip.yqie.com/proxyhttps/
//高匿 http://ip.yqie.com/proxygaoni/
//http http://ip.yqie.com/proxypuni/
type _IPYQie struct {
	SiteName string
	progress int //0初始化 1爬取中 2爬取完成
	maxPage  int64
}

var ipYQie *_IPYQie
var onceYQie sync.Once

func NewIPYQieInstance(opts ...func(o *spiderOption)) *_IPYQie {
	if ipYQie == nil {
		onceYQie.Do(func() {
			opt := &spiderOption{maxPage:6}
			for _,fn := range opts {
				fn(opt)
			}
			ipYQie = &_IPYQie{SiteName: "YQie", progress: 0, maxPage: opt.maxPage}
		})
	}
	return ipYQie
}

func (i *_IPYQie) CrawlData() {
	requestsURL := []requestURL{
		{URL: "http://ip.yqie.com/proxyhttps", Type: "HTTPS"},
		{URL: "http://ip.yqie.com/proxygaoni", Type: "高匿"},
		{URL: "http://ip.yqie.com/proxypuni", Type: "HTTP"},
	}
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	wg := &sync.WaitGroup{}
	for _, v := range requestsURL {
		var dom *goquery.Document
		if dom = getDOMByURL(v.URL); dom == nil {
			utils.Logger().Debugf("IPYQie CrawlData getDOMByURL url:[%v] dom is empty", v.URL)
			continue
		}
		divs := dom.Find("div.divcenter").Find("div")
		if divs.Length() < 4 {
			continue
		}
		div := divs.Eq(3)
		//解析IP数据
		go i.Parser(dom, map[string]interface{}{"url": v.URL})
		if pagenation := div.ChildrenFiltered("div"); pagenation != nil {
			a := pagenation.Find("a")
			if a == nil || a.Length() <= 1 {
				wg.Done()
				continue
			}
			href, exist := a.Eq(1).Attr("href")
			if !exist {
				wg.Done()
				continue
			}
			r, _ := regexp.Compile(`[+]{0,1}(\d+)`)
			page, _ := strconv.Atoi(r.FindString(href))
			utils.Logger().Infof("IPYQie type:[%v],url:[%v],totalPage:[%v]", v.Type, v.URL, page)

			if page > int(i.maxPage) {
				page = int(i.maxPage)
			}
			for j := 2; j <= page; j++ {
				wg.Add(1)
				go func(u string) {
					defer wg.Done()
					if d := getDOMByURL(u); d != nil {
						(*_IPYQie).Parser(&_IPYQie{}, d, map[string]interface{}{"url": u})
					}
				}(fmt.Sprintf("%v/index_%v.htm", v.URL, j))
			}
			wg.Wait()
		}
	}
}
func (ii *_IPYQie) Parser(dom *goquery.Document, params map[string]interface{}) {
	trs := dom.Find("#GridViewOrder").Find("tbody").Find("tr")
	length := trs.Length()
	reg, _ := regexp.Compile(`[document.write(window.atob(\"](\w+\={0,3})[\"));]`)
	for i := 1; i < length; i++ {
		tds := trs.Eq(i).Find("td")
		if tds == nil {
			continue
		}
		ip := tds.Eq(1).Text()
		if f := reg.FindString(ip); f != "" {
			ipByte, e := base64.StdEncoding.DecodeString(strings.Replace(f, "\"", "", -1))
			if e != nil {
				utils.Logger().Errorf("IPYQie Parser base64 decode error: %v,origin data: %v", e, f)
				continue
			}
			ip = string(ipByte)
		}

		port := tds.Eq(2).Text()
		isHTTPs := 0
		if strings.Contains(tds.Eq(4).Text(), "HTTPS") {
			isHTTPs = 1
		}
		local := tds.Eq(3).Text()
		ProxiesChannel <- constructProxy(ip,port,isHTTPs,0,local,"未知",ii.SiteName)
	}
}

func (i *_IPYQie) Progress() int {
	return i.progress
}
func(i *_IPYQie) ProgressName() string{
	return i.SiteName
}