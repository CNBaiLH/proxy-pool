/**
* @Author: Lanhai Bai
* @Date: 2021/8/26 11:04
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

type _IPProxy3366 struct {
	progress int
	SiteName string
	maxPage  int64
}

var ipProxy3366 *_IPProxy3366
var onceProxy3366 sync.Once

func NewIPProxy3366Instance(opts ...func(o *spiderOption)) *_IPProxy3366 {
	if ipProxy3366 == nil {
		onceProxy3366.Do(func() {
			opt := &spiderOption{maxPage: 6}
			for _, fn := range opts {
				fn(opt)
			}
			ipProxy3366 = &_IPProxy3366{progress: 0, SiteName: "齐云代理", maxPage: opt.maxPage}
		})
	}
	return ipProxy3366
}

func (i *_IPProxy3366) CrawlData() {
	defer func() { i.progress = PROGRESS_DONE }()
	i.progress = PROGRESS_DOING
	wg := &sync.WaitGroup{}
	var pageCount int64 = 1
	var crawldata func(u string)
	crawldata = func(u string) {
		defer wg.Done()
		pageCount = atomic.LoadInt64(&pageCount)
		if i.maxPage > 0 && pageCount > i.maxPage {
			utils.Logger().Infof("IPProxy3366 CrawlData Reaches the maximum number of pages: [%v]", pageCount)
			return
		}
		atomic.AddInt64(&pageCount, 1)

		var dom *goquery.Document
		if dom = getDOMByURL(u); dom == nil {
			utils.Logger().Debugf("IPProxy3366 DrawlData dom is empty,url:[%v]", u)
			return
		}
		tbody := dom.Find("div.container").Find("table").Find("tbody")
		if tbody == nil {
			return
		}
		trs := tbody.Find("tr")
		if trs.Length() == 0 {
			return
		}
		i.Paraser(dom, nil)
		pagenation := dom.Find("ul.pagination").Find("li")
		if attr, exsit := pagenation.Last().ChildrenFiltered("a[aria-label=Next]").Attr("href"); exsit {
			wg.Add(1)
			go crawldata("https://proxy.ip3366.net/free/" + attr)
		}
	}
	wg.Add(1)
	go crawldata("https://proxy.ip3366.net/free")
	wg.Wait()
}

func (ii *_IPProxy3366) Paraser(dom *goquery.Document, params map[string]interface{}) {
	utils.Logger().Debug("_IPProxy3366 Paraser...")
	tbody := dom.Find("div.container").Find("table").Find("tbody")
	trs := tbody.Find("tr")
	length := trs.Length()
	for i := 0; i < length; i++ {
		tds := trs.Eq(i).Find("td")
		var isHTTPs = 0
		if strings.Contains(tds.Eq(3).Text(), "HTTPS") {
			isHTTPs = 1
		}
		ip := tds.Eq(0).Text()
		port := tds.Eq(1).Text()
		local := tds.Eq(4).Text()
		ProxiesChannel <- constructProxy(ip, port, isHTTPs, 0, local, "未知", ii.SiteName)
	}
}

func (i *_IPProxy3366) Progress() int {
	return i.progress
}
func (i *_IPProxy3366) ProgressName() string {
	return i.SiteName
}
