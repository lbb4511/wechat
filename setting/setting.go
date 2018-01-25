// Package setting provides ...
package setting

import (
	"io/ioutil"
	"net/http"
	"time"

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

func WriteLog(r *http.Request, t time.Time, match string, pattern string) {

	if PROD == Conf.RunMode {

		d := time.Now().Sub(t)

		log.Printf("[ACCESS] | % -10s | % -40s | % -16s | % -10s | % -40s |\n", r.Method, r.URL.Path, d.String(), match, pattern)
	} else if DEV == Conf.RunMode {

		log.Print(DEV)

	} else {
		log.Fatal("err")
	}
}
