package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type Pager struct {
	Count       int64 `json:"count"`
	TotalPage   int64 `json:"total_page"`
	CurrentPage int `json:"page"`
	PerPage     int `json:"per_page"`
}
const pageKey = "page" //页码请求参数
const pageLimit = "per_page" //每页显示数量

func newPager(page,perPage int) *Pager{
	return &Pager{
		Count:       0,
		TotalPage:   0,
		CurrentPage: page,
		PerPage:     perPage,
	}
}

//@function 初始化分页对象
func InitPage(ctx *gin.Context) *Pager {
	//p := &Pager{}
	//当前页
	pageStr := ctx.DefaultQuery(pageKey, "1")
	//每页显示多少项
	perPageStr := ctx.DefaultQuery(pageLimit, "10")
	var page int
	var perPage float64
	page, _ = strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	perPage, _ = strconv.ParseFloat(perPageStr, 64)
	if perPage <= 0 {
		perPage = 10
	}
	return newPager(page,int(perPage))
}

//@function 设置分页总条数
func (p *Pager) SetCount(count int64) {
	p.Count = count
	totalPage := math.Ceil(float64(count) / float64(p.PerPage))
	p.TotalPage = int64(totalPage)
}

//@function 取得当前设置的limit
//@return Limit(10,0) => limit 10 offset 0
func (p *Pager) GetLimit() (int, int) {
	page := p.CurrentPage
	perPage := p.PerPage
	return perPage, (page - 1) * perPage
}

//@function 返回limit的拼接字符串
//@return 格式 limit 条数 offset start
func (p *Pager) GetLimitString() (limitString string) {
	return fmt.Sprintf("%d offset %d", p.PerPage, (p.CurrentPage-1)*p.PerPage )
}

//@function 返回当前页对应的offset
//@return int 偏移量 第一页偏移0条  那第二页就偏移10条
func (p *Pager) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}
