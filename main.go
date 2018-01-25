package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/lbb4511/wechat/controller"
	"github.com/lbb4511/wechat/setting"
)

type httpHandler struct{}

type WebController struct {
	Function func(http.ResponseWriter, *http.Request)
	Method   string
	Pattern  string
}

var mux []WebController

func init() {
	mux = append(mux, WebController{controller.Post, "POST", "^/wechat"})
	mux = append(mux, WebController{controller.Get, "GET", "^/wechat"})
}

func main() {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Conf.Port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", setting.Conf.Port))
	log.Fatal(server.ListenAndServe())
}

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	for _, webController := range mux {

		if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m {

			if r.Method == webController.Method {

				webController.Function(w, r)

				go setting.WriteLog(r, t, "match", webController.Pattern)

				return
			}
		}
	}

	go setting.WriteLog(r, t, "unmatch", "")

	io.WriteString(w, "")
	return
}
