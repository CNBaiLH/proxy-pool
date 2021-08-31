/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 11:08
* @Description:
 */
package service

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"proxy-pool/model"
	"proxy-pool/utils"
	"strings"
	"time"
)

const (
	PROGRESS_INIT  = 0 //初始化
	PROGRESS_DOING = 1 //爬取中
	PROGRESS_DONE  = 2 //完成
)

var ProxiesChannel chan model.Proxy

func init() {
	ProxiesChannel = make(chan model.Proxy, 100)
}

type spiderOption struct {
	maxPage int64
}

func WithSpiderRequestMaxPage(maxPage int64) func(o *spiderOption) {
	return func(o *spiderOption) {
		if maxPage >= 1 {
			o.maxPage = maxPage
		} else {
			o.maxPage = 6
		}
	}
}

//爬取数据
type CrawlDataInterfacer interface {
	CrawlData()
}

//解析数据
type ParserDataInterfacer interface {
	Parser(dom *goquery.Document, params map[string]interface{})
}

//进度查询
type ProgressInterfacer interface {
	Progress() int
	ProgressName() string
}

type RequestDOMOption struct {
	duration time.Duration     //当前请求延迟N秒后执行，如果是默认值 -1 则随机一个时间， 如果为0则立即执行
	headers  map[string]string //请求头
	timeout  time.Duration     //请求超时时间
}

//设置请求延迟时间
func WithHTTPDuration(duration time.Duration) func(o *RequestDOMOption) {
	return func(o *RequestDOMOption) {
		o.duration = duration
	}
}

//设置请求头
func WithHTTPHeader(headers map[string]string) func(o *RequestDOMOption) {
	return func(o *RequestDOMOption) {
		o.headers = headers
	}
}

//设置超时时间
func WithHTTPTimeout(timeout time.Duration) func(o *RequestDOMOption) {
	return func(o *RequestDOMOption) {
		o.timeout = timeout
	}
}

func getDOMByURL(cURL string, opts ...func(option *RequestDOMOption)) *goquery.Document {
	opt := &RequestDOMOption{
		duration: -1, //随机5秒内的时间再去访问
		headers:  nil,
		timeout:  time.Duration(time.Second * 20),
	}
	for _, fn := range opts {
		fn(opt)
	}
	if opt.duration == -1 {
		opt.duration = time.Duration(rand.Intn(5)) * time.Second
	}
	request, _ := http.NewRequest("GET", cURL, nil)
	client := &http.Client{
		Transport: &http.Transport{
		},
		Timeout: opt.timeout,
	}
	for k, v := range opt.headers {
		request.Header.Set(k, v)
	}
	utils.Logger().Infof("getDOMByURL Request URL: %v", cURL)
	rsp, err := client.Do(request)
	if err != nil {
		utils.Logger().Errorf("getDOMByURL error: %v,url: %v", err, cURL)
		return nil
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		utils.Logger().Debugf("getDOMByURL response code :[%v],body :%s,url:[%v]", rsp.StatusCode, body, cURL)
		return nil
	}
	dom, _ := goquery.NewDocumentFromReader(ioutil.NopCloser(strings.NewReader(string(body))))
	return dom
}

func replaceSpaceAndOtherMark(ll string) string {
	ll = strings.Replace(ll, "\t", "", -1)
	ll = strings.Replace(ll, "\n", "", -1)
	ll = strings.Replace(ll, " ", "", -1)
	return ll
}

type requestURL struct {
	URL  string
	Type string
}

//检查代理是否可用
func IsLive(proxyURL string) bool {
	testURL := "http://icanhazip.com"
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		utils.Logger().Errorf("IsLive Parse error: %v", err)
		return false
	}
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(8),
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	res, err := httpClient.Get(testURL)
	if err != nil {
		utils.Logger().Errorf("IsLive httpClient.Get [%v] error: %v", proxyURL, err)
		return false
	}
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		utils.Logger().Debugf("IsLive response code:[%v], body:[%s]", res.StatusCode, body)
		return false
	}
	return true
}

func constructProxy(ip string, port string, isHTTPs int, isPost int, local string, isp string, source string) model.Proxy {
	var status int = 1 //刚入库的数据默认可用

	p := model.Proxy{
		Host:         ip,
		Port:         port,
		Local:        local,
		SupportHttps: isHTTPs,
		SupportPost:  isPost,
		Isp:          isp,
		Created:      time.Now().Unix(),
		CheckTime:    time.Now().Unix(),
		Status:       status,
		Source:       source,
	}
	utils.Logger().Infof("constructProxy proxy info: %+v", p)
	return p
}
