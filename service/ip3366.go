/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 15:23
* @Description:
 */
package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"proxy-pool/utils"
	"strconv"
	"strings"
	"sync"
)

type _IP3366 struct {
	progress int
	SiteName string
	maxPage  int64
}

var ip3366 *_IP3366
var onceIP3366 sync.Once

func NewIP3366Instance(opts ...func(o *spiderOption)) *_IP3366 {
	if ip3366 == nil {
		onceIP3366.Do(func() {
			opt := &spiderOption{maxPage:6}
			for _,fn := range opts {
				fn(opt)
			}
			ip3366 = &_IP3366{progress: 0, SiteName: "ip3366云代理", maxPage: opt.maxPage}
		})
	}
	return ip3366
}

func (i *_IP3366) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	requests := []requestURL{
		{URL: "http://www.ip3366.net/free/?stype=1", Type: "高匿"},
		{URL: "http://www.ip3366.net/free/?stype=2", Type: "普通"},
	}
	wg := &sync.WaitGroup{}
	for _, v := range requests {
		var dom *goquery.Document
		if dom = getDOMByURL(v.URL); dom == nil {
			utils.Logger().Debugf("IP3366 CrawlData getDOMByURL url:[%v] dom is empty", v.URL)
			continue
		}
		pagenationText := dom.Find("#listnav").Find("ul").ChildrenFiltered("strong").Text()
		pagenation := strings.Split(pagenationText, "/")
		var page = 1
		if len(pagenation) == 2 {
			page, _ = strconv.Atoi(pagenation[1])
		}
		if page > int(i.maxPage) {
			page = int(i.maxPage)
		}
		for j := 2; j <= page; j++ {
			wg.Add(1)
			go func(baseURL string, p int) {
				defer wg.Done()
				if d := getDOMByURL(fmt.Sprintf("%v&page=%v", baseURL, p)); d != nil {
					i.Parser(d, nil)
				}
			}(v.URL, j)
		}
		i.Parser(dom, nil)
	}
	wg.Wait()
}

func (ii *_IP3366) Parser(dom *goquery.Document, params map[string]interface{}) {
	tbody := dom.Find("#list").Find("tbody")
	if tbody == nil {
		return
	}
	trs := tbody.Find("tr")
	length := trs.Length()
	for i := 0; i < length; i++ {
		tds := trs.Eq(i).Find("td")
		ip := tds.Eq(0).Text()
		port := tds.Eq(1).Text()
		isPost := 0
		if strings.Contains(tds.Eq(4).Text(), "POST") {
			isPost = 1
		}
		isHTTPs := 0
		if tds.Eq(3).Text() == "HTTPS" {
			isHTTPs = 1
		}
		local := tds.Eq(4).Text()
		dec := mahonia.NewDecoder("gbk")
		local = dec.ConvertString(local)
		ProxiesChannel <- constructProxy(ip,port,isHTTPs,isPost,local,"未知",ii.SiteName)
	}
}

func (i *_IP3366) Progress() int {
	return i.progress
}
func(i *_IP3366) ProgressName() string{
	return i.SiteName
}