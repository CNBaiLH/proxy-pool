/**
* @Author: Lanhai Bai
* @Date: 2021/8/26 16:30
* @Description:
 */
package model

import (
	"proxy-pool/utils"
	"time"
	"xorm.io/builder"
)

type Proxy struct {
	Id           int64  `json:"id" xorm:"not null autoincr pk INT(10)"`
	Host         string `json:"host" xorm:"unique(idx_host_port) index(idx_host_port_status_https_post) not null default '' VARCHAR(18)"`
	Port         string `json:"port" xorm:"unique(idx_host_port) index(idx_host_port_status_https_post) not null default '0' VARCHAR(8)"` //端口
	Local        string `json:"local" xorm:"not null default '' VARCHAR(60)"`                                                             //位置
	SupportHttps int    `json:"support_https" xorm:"index(idx_host_port_status_https_post) not null default 0 comment('是否支持HTTPS 1支持 0不支持') TINYINT(2)"`
	SupportPost  int    `json:"support_post" xorm:"index(idx_host_port_status_https_post) not null default 0 comment('是否支持POST 1支持 0不支持') TINYINT(2)"`
	Isp          string `json:"isp" xorm:"not null default '' comment('运营商') VARCHAR(30)"`
	Created      int64  `json:"created" xorm:"'created' not null default 0 INT(10)"` //入库时间
	CheckTime    int64  `json:"check_time" xorm:"not null default 0 INT(10)"`        //检测时间
	Status       int    `json:"status" xorm:"index(idx_host_port_status_https_post) not null default 0 comment('存活状态 1存活 0未存活') TINYINT(2)"`
	Source       string `json:"-" xorm:"not null default '未知' comment('数据源') VARCHAR(30)"`
}

func MutilSaveProxies(proxies []Proxy) bool {
	sql := `INSERT IGNORE INTO proxy (host, port, status,created,check_time,local,isp,support_https,support_post,source) VALUES `
	values := []interface{}{sql}
	for _, proxy := range proxies {
		values = append(values, proxy.Host, proxy.Port, proxy.Status, proxy.Created, proxy.CheckTime, proxy.Local, proxy.Isp, proxy.SupportHttps, proxy.SupportPost, proxy.Source)
		sql += "(?,?,?,?,?,?,?,?,?,?),"
	}
	sql = sql[0:len(sql)-1] + ";"
	values[0] = sql
	utils.Logger().Infof("保存数据的sql为:%v", sql)
	result, err := Engine().Exec(values...)
	if err != nil {
		utils.Logger().Errorf("MutilSaveProxies insert into proxy error: %v", err)
	}
	affect, _ := result.RowsAffected()
	return affect > 0
}
func SaveProxy(p *Proxy) bool {
	sql := `INSERT INTO proxy
(host, port, status,created,check_time,local,isp,support_https,support_post)
SELECT ?,?,?,?,?,?,?,?,?
FROM DUAL
WHERE not exists (select * from proxy t
where t.host = ? and t.port=?);`
	result, err := Engine().Exec(sql, p.Host, p.Port, p.Status, p.Created, p.CheckTime, p.Local, p.Isp, p.SupportHttps, p.SupportPost, p.Host, p.Port)
	if err != nil {
		utils.Logger().Errorf("SaveProxy insert into proxy error: %v", err)
	}
	affect, _ := result.RowsAffected()
	return affect > 0
}

func QueryPoxiesByCond(cond builder.Cond, limit int, offset int, orderBy string) []Proxy {
	proxies := []Proxy{}
	if err := Engine().Table("proxy").Where(cond).Limit(limit, offset).OrderBy(orderBy).Find(&proxies); err != nil {
		utils.Logger().Errorf("QueryPoxiesByCond query proxy error: %v", err)
	}
	return proxies
}

func QueryPoxiesAndCountByCond(cond builder.Cond, limit int, offset int, orderBy string) ([]Proxy, int64) {
	proxies := []Proxy{}
	var count int64
	var err error
	if count, err = Engine().Table("proxy").Where(cond).Limit(limit, offset).OrderBy(orderBy).FindAndCount(&proxies); err != nil {
		utils.Logger().Errorf("QueryPoxiesAndCountByCond query proxy error: %v", err)
	}
	return proxies, count
}

//更新ip的状态
func UpdateProxyStatusById(id int64, status int) {
	utils.Logger().Debugf("UpdateProxyStatusById update proxy data, id:[%v]", id)
	if _, err := Engine().Table("proxy").Where("id=?", id).Cols("status", "check_time").Update(map[string]interface{}{"status": status, "check_time": time.Now().Unix()}); err != nil {
		utils.Logger().Errorf("UpdateProxyStatusById update proxy error: %v", err)
	}
}

//删除指定ip
func DeleteProxyById(id int64) {
	if _, err := Engine().Table("proxy").Where("id=?", id).Delete(new(Proxy)); err != nil {
		utils.Logger().Errorf("DeleteProxyById delete proxy error: %v,id:[%v]", err, id)
	}
}

//查询数据条数
func CountProxy() (count int64) {
	count, err := Engine().Table("proxy").Where("id>0").Count(new(Proxy))
	if err != nil {
		utils.Logger().Errorf("CountProxy query proxy count error: %v", err)
	}
	return count
}
