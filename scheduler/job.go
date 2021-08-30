/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 10:10
* @Description:
 */
package scheduler

import (
	"github.com/robfig/cron"
	"proxy-pool/model"
	"proxy-pool/service"
	"proxy-pool/utils"
	"time"
	"xorm.io/builder"
)

func init() {
	utils.Logger().Infof("job init...")
	cr := cron.New()
	//每天凌晨1点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 1 * * ?", func() {
		CollectProxyData()
	})

	//每天早上5点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 5 * * ?", func() {
		CollectProxyData()
	})

	//每天早上10点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 10 * * ?", func() {
		CollectProxyData()
	})

	//每天14点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 14 * * ?", func() {
		CollectProxyData()
	})

	//每天18点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 18 * * ?", func() {
		CollectProxyData()
	})


	//每天22点自动爬取不同站点的IP列表
	cr.AddFunc("0 0 22 * * ?", func() {
		CollectProxyData()
	})

	//每过5分钟检测一下存活状态
	cr.AddFunc("0 */5 * * * ?", func() {
		CheckProxyStatusAndUpdateProxyInfo()
	})
	cr.Start()
	//系统启动执行数据检测，一条数据没有则立马进行数据爬取
	go ExecOnce()
	go SynchronizeToTheDatabase()
}

//系统初始化执行
//如果没有数据则立马进行爬取
func ExecOnce() {
	if model.CountProxy() == 0 {
		utils.Logger().Info("ExecOnce proxy count 0,run CollectProxyData at now...")
		CollectProxyData()
	}
}

//同步数据到数据库
func SynchronizeToTheDatabase() {
	utils.Logger().Infof("SynchronizeToTheDatabase begin...")
	//每1s执行数据同步操作
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	proxies := []model.Proxy{}
	for {
		select {
		case <-ticker.C:
			if len(proxies) > 0 {
				utils.Logger().Debugf("执行批量保存...")
				model.MutilSaveProxies(proxies)
				utils.Logger().Debugf("批量保存执行完毕...")
				proxies = []model.Proxy{}
			}
		case p := <-service.ProxiesChannel:
			proxies = append(proxies, p)
		default:
			<-time.After(time.Millisecond * 200)
		}
	}
}

//收集所有数据
func CollectProxyData() {
	utils.Logger().Infof("CollectProxyData begin...")
	defer utils.Logger().Infof("CollectProxyData end...")
	var proxySourceList = []interface{}{
		service.NewIP3366Instance(),
		service.NewIPKuaiDailiInstance(),
		service.NewIPProxy3366Instance(),
		service.NewIPXilaInstance(service.WithSpiderRequestMaxPage(3)),
		service.NewIPYQieInstance(),
		service.NewIP89Instance(),
		service.NewIPIHuanInstance(),
	}

	for _, v := range proxySourceList {
		if p, ok := v.(service.ProgressInterfacer); ok {
			if p.Progress() == 1 {
				utils.Logger().Debugf("Collection in progress:[%v]", p.ProgressName())
				continue
			}
		}
		go v.(service.CrawlDataInterfacer).CrawlData()
	}
}

//检查代理存活情况并更新数据
func CheckProxyStatusAndUpdateProxyInfo() {
	utils.Logger().Infof("CheckProxyStatusAndUpdateProxyInfo begin...")
	defer utils.Logger().Infof("CheckProxyStatusAndUpdateProxyInfo end...")

	currentTime := time.Now().Unix()
	cond := builder.NewCond().And(builder.Eq{"status": 0}).And(builder.Lte{"check_time": currentTime - 30*60})
	proxies := model.QueryPoxiesByCond(cond, 400, 0, "created ASC")
	for _, p := range proxies {
		go checkProxyStatusAndUpdate(p)
	}
}

//检查状态 并更新
func checkProxyStatusAndUpdate(p model.Proxy) {
	var status = 0
	if service.IsLive("http://" + p.Host + ":" + p.Port) {
		status = 1
	}
	//15天之前入库的IP  且现在是未存活状态，则删除掉
	if p.Created < time.Now().Unix()-15*86400 && status == 0 {
		model.DeleteProxyById(p.Id)
		return
	}
	model.UpdateProxyStatusById(p.Id, status)
}
