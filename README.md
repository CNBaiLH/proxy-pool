# proxy-pool

<p align="center">
  <img src="https://img.shields.io/badge/golang-v1.15.7-green" />
</p>

**项目基于 `go` + `goquery` + `gin` + `xorm` 实现的 `IP池` 服务，每日多个时段进行IP池数据扩充，定时进行存活检测。**
**采用mysql作为数据持久化存储**


<u>`写得不好，希望大佬帮忙指正，万分感谢`</u>

### 数据源
- [89代理](https://www.89ip.cn)
- [小幻代理](https://ip.ihuan.me)
- [ip3366云代理](http://www.ip3366.net)
- [快代理](https://www.kuaidaili.com/free)
- [齐云代理](https://proxy.ip3366.net/free)
- [西拉代理](http://www.xiladaili.com)
- [YQie代理](http://ip.yqie.com)

## TODO
- [ ]  扩充IP数据源
- [ ]  爬虫优化

## 🚀 如何运行

### 1.克隆
```
git clone https://github.com/CNBaiLH/proxy-pool.git
```

### 2.进入项目目录
```
cd proxy-pool
```

## 3.快速运行
> 3.1 Docker方式运行
```
docker-compose -f ./docker-compose.yaml up
```

> 3.2 普通方式运行
```
go run main.go
```


## 📖API目录
### 查询满足条件的IP列表
**GET**  **`http://host:port/v1/ips`**

| 参数        | 是否必须   |  参考值  |
| --------   | -----:  | :----:  |
| status     |  否     |   0(未存活)或1(存活)  |
| https      |   否    |  0(不支持)或1(支持)  |
| post       |    否   |  0(不支持)或1(支持)  |
| page       |    否   | 当前的页码        |
| per_page   |    否   | 每页显示条数    |


> 1. 请求示例 

```1.1 查询所有的IP列表```

***curl `http://127.0.0.1/v1/ips?page=1&per_page=1`***

```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "list": [
            {
                "id": 310122,
                "host": "95.210.251.29",
                "port": "53281",
                "local": "意大利 意大利 ",
                "support_https": 0,
                "support_post": 1,
                "isp": "skylogic.it",
                "created": 1630028464,
                "check_time": 1630035605,
                "status": 1
            }
        ],
        "meta": {
            "count": 2251,
            "total_page": 2251,
            "page": 1,
            "per_page": 1
        }
    }
}
```


```1.2 查询存活的IP列表```

***curl `http://127.0.0.1/v1/ips?page=1&per_page=1&status=1`***
```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "list": [
            {
                "id": 310122,
                "host": "95.210.251.29",
                "port": "53281",
                "local": "意大利 意大利 ",
                "support_https": 0,
                "support_post": 1,
                "isp": "skylogic.it",
                "created": 1630028464,
                "check_time": 1630035605,
                "status": 1
            }
        ],
        "meta": {
            "count": 393,
            "total_page": 393,
            "page": 1,
            "per_page": 1
        }
    }
}
```



## 项目目录结构

proxy-pool
> handlers -- 路由控制器
>
> logs --运行日志
>
> model --数据库模型
>
> router --路由
>
> scheduler --后台任务
>
> service --实现各个代理网站的爬取
>
> utils --工具类



## 👤作者

[菜鸟ITWorker](https://www.cnblogs.com/5566blh/)