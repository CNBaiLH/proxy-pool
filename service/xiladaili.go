/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 15:53
* @Description:
 */
package service

import (
	"github.com/PuerkitoBio/goquery"
	"proxy-pool/utils"
	"strings"
	"sync"
	"sync/atomic"
)

//http://www.xiladaili.com/gaoni 高匿
//http://www.xiladaili.com/http http
//http://www.xiladaili.com/https https代理
type _IPXila struct {
	progress int
	SiteName string
	maxPage  int64
}

var ipXila *_IPXila
var onceIPXila sync.Once

func NewIPXilaInstance(opts ...func(o *spiderOption)) *_IPXila {
	if ipXila == nil {
		onceIPXila.Do(func() {
			opt := &spiderOption{maxPage:6}
			for _,fn := range opts {
				fn(opt)
			}
			ipXila = &_IPXila{progress: 0, SiteName: "西拉代理", maxPage: opt.maxPage}
		})
	}
	return ipXila
}

func (i *_IPXila) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	requests := []requestURL{
		{URL: "http://www.xiladaili.com/gaoni/", Type: "高匿"},
		{URL: "http://www.xiladaili.com/http/", Type: "HTTP"},
		{URL: "http://www.xiladaili.com/https/", Type: "HTTPS"},
	}
	var pageDict = map[string]*int64{}
	wg := &sync.WaitGroup{}
	var crawldata func(u string,baseURL string, isFirst bool)
	crawldata = func(u string,baseURL string, isFirst bool) {
		defer wg.Done()
		pageCount := atomic.LoadInt64(pageDict[baseURL])
		if i.maxPage > 0 && pageCount > i.maxPage {
			utils.Logger().Infof("IPXila CrawlData Reaches the maximum number of pages: [%v]", pageCount)
			return
		}
		atomic.AddInt64(pageDict[baseURL], 1)

		var dom *goquery.Document
		if dom = getDOMByURL(u,WithHTTPHeader(map[string]string{"Cookie":"Hm_lvt_9bfa8deaeafc6083c5e4683d7892f23d=1629689060,1629787777,1629941587,1630375400; Hm_lpvt_9bfa8deaeafc6083c5e4683d7892f23d=1630375457","Upgrade-Insecure-Requests":"1","Proxy-Connection":"keep-alive","User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"})); dom == nil {
			utils.Logger().Debugf("IPXila CrawlData getDOMByURL url:[%v] dom is empty", u)
			return
		}
		lis := dom.Find("ul.pagination").Find("li.page-item")
		if lis == nil {
			return
		}

		length := lis.Length()
		nextPageLi := lis.Eq(length - 1)
		trs := dom.Find("table.fl-table").Find("tbody").Find("tr")
		if nextPageLi != nil && trs.Length() != 0 {
			wg.Add(1)
			href, _ := nextPageLi.Children().Attr("href")
			go crawldata("http://www.xiladaili.com"+href,baseURL, false)
			return
		}
		i.Parser(dom, nil)
	}

	for _, v := range requests {
		var pInt int64 = 1
		pageDict[v.URL] = &pInt
		wg.Add(1)
		go crawldata(v.URL,v.URL, true)
	}
	wg.Wait()
}

func (ii *_IPXila) Parser(dom *goquery.Document, params map[string]interface{}) {
	utils.Logger().Debug("_IPXila Parser...")
	tbody := dom.Find("table.fl-table").Find("tbody")
	if tbody == nil {
		return
	}
	trs := tbody.Find("tr")
	length := trs.Length()
	for i := 0; i < length; i++ {
		tds := trs.Eq(i).Find("td")
		addr := strings.Split(tds.Eq(0).Text(), ":")
		isHTTPs := 0
		if strings.Contains(tds.Eq(1).Text(), "HTTPS") {
			isHTTPs = 1
		}
		var ip, port string
		if len(addr) == 2 {
			ip = addr[0]
			port = addr[1]
		}
		local := tds.Eq(3).Text()

		ProxiesChannel <- constructProxy(ip,port,isHTTPs,0,local,"未知",ii.SiteName)
	}
}

func (i *_IPXila) Progress() int {
	return i.progress
}
func(i *_IPXila) ProgressName() string{
	return i.SiteName
}
