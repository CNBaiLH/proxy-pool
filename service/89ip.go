/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 14:35
* @Description:
 */
package service

import (
	"github.com/PuerkitoBio/goquery"
	"proxy-pool/utils"
	"sync"
	"sync/atomic"
)

type _IP89 struct {
	SiteName string
	progress int   //0初始化 1爬取中 2爬取完成
	maxPage  int64 //最大爬取页数
}

var ip89 *_IP89
var onceIP89 sync.Once

func NewIP89Instance(opts ...func(o *spiderOption)) *_IP89 {
	if ip89 == nil {
		onceIP89.Do(func() {
			opt := &spiderOption{maxPage: 6}
			for _, fn := range opts {
				fn(opt)
			}
			ip89 = &_IP89{
				SiteName: "89代理",
				progress: 0,
				maxPage:  opt.maxPage,
			}
		})
	}
	return ip89
}

//https://www.89ip.cn
func (i *_IP89) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	var crawldata func(u string, isFirst bool)
	var pageCount int64 = 1
	wg := &sync.WaitGroup{}
	crawldata = func(u string, isFirst bool) {
		defer wg.Done()
		pageCount = atomic.LoadInt64(&pageCount)
		if i.maxPage > 0 && pageCount > i.maxPage {
			utils.Logger().Infof("IP89 CrawlData Reaches the maximum number of pages: [%v]", pageCount)
			return
		}
		atomic.AddInt64(&pageCount, 1)

		var dom *goquery.Document
		if dom = getDOMByURL(u); dom == nil {
			utils.Logger().Debugf("IP89 CrawlData getDOMByURL url:[%v] dom is empty", u)
			return
		}
		nextPage := dom.Find("a.layui-laypage-next")
		tbody := dom.Find(".layui-table").Find("tbody")
		if tbody == nil {
			return
		}
		trs := tbody.Find("tr")
		if trs.Length() == 0 {
			return
		}
		if nextPage != nil {
			wg.Add(1)
			href, _ := nextPage.Attr("href")
			go crawldata("https://www.89ip.cn/"+href, false)
		}
		i.Parser(dom, nil)
	}
	wg.Add(1)
	go crawldata("https://www.89ip.cn", true)
	wg.Wait()
}

func (i *_IP89) Parser(dom *goquery.Document, params map[string]interface{}) {
	utils.Logger().Debug("_IP89 Parser...")
	tbody := dom.Find(".layui-table").Find("tbody")
	trs := tbody.Find("tr")
	length := trs.Length()
	for j := 0; j < length; j++ {
		tds := trs.Eq(j).Find("td")
		ip := replaceSpaceAndOtherMark(tds.Eq(0).Text())
		port := replaceSpaceAndOtherMark(tds.Eq(1).Text())
		isHTTPs := 0
		isp := replaceSpaceAndOtherMark(tds.Eq(3).Text())
		local := replaceSpaceAndOtherMark(tds.Eq(2).Text())
		ProxiesChannel <- constructProxy(ip, port, isHTTPs, 0, local, isp, i.SiteName)
	}
}

func (i *_IP89) Progress() int {
	return i.progress
}

func (i *_IP89) ProgressName() string {
	return i.SiteName
}
