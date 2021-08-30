package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"proxy-pool/utils"
	"time"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func init() {
	utils.Logger().Infof("model init...")
	engine = new(xorm.Engine)
	config := utils.Configure()
	initDB(config.Database.Host, config.Database.Port, config.Database.Password, config.Database.Username, config.Database.DBName, )
	if err := engine.Sync2(new(Proxy)); err != nil {
		utils.Logger().Errorf("model init sync2 error:%v", err)
	}
}

func initDB(dbHost, dbPort, dbPassWD, dbUserName, dbName string, ) {
	utils.Logger().Infof("db init...")
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", dbUserName, dbPassWD, dbHost, dbPort, dbName, "utf8")
	var err error
	engine, err = xorm.NewEngine("mysql", dbSource)
	if err != nil {
		panic("InitDB new engine error:" + err.Error())
	}
	engine.ShowSQL(true)
	engine.SetMaxOpenConns(20)
	engine.SetConnMaxLifetime(time.Duration(10) * time.Second)
}

func Engine() *xorm.Engine {
	return engine
}
