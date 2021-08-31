/**
* @Author: Lanhai Bai
* @Date: 2021/8/26 11:37
* @Description:小幻代理
 */
package service

import (
	"github.com/PuerkitoBio/goquery"
	"proxy-pool/utils"
	"sync"
	"sync/atomic"
)

type _IPIHuan struct {
	SiteName string
	progress int
	maxPage  int64
}

var onceIPIHuan sync.Once
var ipIhuan *_IPIHuan

func NewIPIHuanInstance(opts ...func(o *spiderOption)) *_IPIHuan {
	if ipIhuan == nil {
		onceIPIHuan.Do(func() {
			opt := &spiderOption{maxPage: 6}
			for _, fn := range opts {
				fn(opt)
			}
			ipIhuan = &_IPIHuan{progress: 0, SiteName: "小幻代理", maxPage: opt.maxPage}
		})
	}
	return ipIhuan
}

//https://ip.ihuan.me/
func (i *_IPIHuan) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	var crawldata func(u string)
	wg := &sync.WaitGroup{}
	var pageCount int64 = 1
	crawldata = func(u string) {
		defer wg.Done()
		pageCount = atomic.LoadInt64(&pageCount)
		if i.maxPage > 0 && pageCount > i.maxPage {
			utils.Logger().Infof("IPIHuan CrawlData Reaches the maximum number of pages: [%v]", pageCount)
			return
		}
		atomic.AddInt64(&pageCount, 1)

		var dom *goquery.Document
		if dom = getDOMByURL(u); dom == nil {
			utils.Logger().Debugf("IPIHuan CrawlData getDOMByURL url:[%v] dom is empty", u)
			return
		}
		tr := dom.Find(".table-responsive tbody").Find("tr")
		length := tr.Length()
		if length == 0 {
			return
		}
		i.Parser(dom, nil)
		li := dom.Find("ul.pagination").Find(`li`).Last()
		a := li.ChildrenFiltered("a[aria-label=Next]")
		if a == nil {
			return
		}
		if href, exist := a.Attr("href"); exist {
			wg.Add(1)
			go crawldata("https://ip.ihuan.me" + href)
			return
		}
	}
	wg.Add(1)
	go crawldata("https://ip.ihuan.me")
	wg.Wait()
}

func (ii *_IPIHuan) Parser(dom *goquery.Document, params map[string]interface{}) {
	tr := dom.Find(".table-responsive tbody").Find("tr")
	length := tr.Length()
	for i := 0; i < length; i++ {
		tds := tr.Eq(i).Find("td")
		ip := tds.Eq(0).Text()
		port := tds.Eq(1).Text()
		isSupportPost := 0
		if tds.Eq(5).Text() == "支持" {
			isSupportPost = 1
		}
		isHTTPs := 0
		if tds.Eq(4).Text() == "支持" {
			isHTTPs = 1
		}
		isp := tds.Eq(3).Text()
		local := tds.Eq(2).Text()

		ProxiesChannel <- constructProxy(ip, port, isHTTPs, isSupportPost, local, isp, ii.SiteName)
	}
}

func (i *_IPIHuan) Progress() int {
	return i.progress
}
func (i *_IPIHuan) ProgressName() string {
	return i.SiteName
}
