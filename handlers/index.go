/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 10:22
* @Description:
 */
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-pool/model"
	"proxy-pool/utils"
	"strconv"
	"strings"
	"xorm.io/builder"
)

//所有的代理
func IndexHandler(ctx *gin.Context) {
	pager := utils.InitPage(ctx)
	limit, offset := pager.GetLimit()
	status := strings.Trim(ctx.DefaultQuery("status", "-1"), " ")
	isHTTPs := strings.Trim(ctx.DefaultQuery("https", "-1"), " ")
	isPost := strings.Trim(ctx.DefaultQuery("post", "-1"), " ")
	cond := builder.NewCond()
	if status != "-1" {
		s, _ := strconv.Atoi(status)
		cond = cond.And(builder.Eq{"status": s})
	}

	if isHTTPs != "-1" {
		i, _ := strconv.Atoi(isHTTPs)
		cond = cond.And(builder.Eq{"is_support_https": i})
	}

	if isPost != "-1" {
		i, _ := strconv.Atoi(isPost)
		cond = cond.And(builder.Eq{"is_support_post": i})
	}
	utils.Logger().Debugf("IndexHandler recive a request... page[%v],per_page[%v],status[%v],https[%v],post[%v]",pager.CurrentPage,pager.PerPage,status,isHTTPs,isPost)
	proxy, count := model.QueryPoxiesAndCountByCond(cond, limit, offset, "created DESC")
	pager.SetCount(count)
	ctx.JSON(http.StatusOK, utils.JSONDataWithPage(proxy, pager))
}