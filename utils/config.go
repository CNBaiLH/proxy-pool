/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 9:43
* @Description:
 */
package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var conf = &allConfig{}

type allConfig struct {
	content  []byte `yaml:"-"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	System struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"system"`
}

func init() {
	Logger().Info("config init...")
	initConf()
}

func initConf() {
	configName := "config.dev.yaml"
	if os.Getenv("RUN_ENV") == "DOCKER" {
		configName = "config.yaml"
	}
	f, err := os.OpenFile(GetWebDir()+"/"+configName, os.O_RDONLY, 0666)
	if err != nil {
		panic("initConf error:" + err.Error())
	}

	content, _ := ioutil.ReadAll(f)
	if err = yaml.Unmarshal(content, conf); err != nil {
		panic("initConf yaml.Unmarshal error:" + err.Error())
	}
	conf.content = content
}

func Configure() *allConfig {
	if conf == nil {
		initConf()
	}
	return conf
}
