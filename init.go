// Package setting provides ...
package main

import (
	"io/ioutil"

	"github.com/lbb4511/wechat/log"
	yaml "gopkg.in/yaml.v2"
)

const (
	DEV  = "dev"  // 该模式会输出 debug 等信息
	PROD = "prod" // 该模式用于生产环境
)

var (
	Conf = new(Config)
)

type Config struct {
	RunMode       string // 运行模式
	StaticVersion int    // 当前静态文件版本
	Port          int    // 端口
	Token         string // Token
}

func init() {
	// 初始化配置
	data, err := ioutil.ReadFile("conf/wechat.ini")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, Conf)
	if err != nil {
		log.Fatal(err)
	}
}
