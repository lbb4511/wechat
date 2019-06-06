package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/lbb4511/wechat/service"
)

type httpHandler struct{}

type WebController struct {
	Function func(http.ResponseWriter, *http.Request)
	Method   string
	Pattern  string
}

var mux []WebController

func main() {

	mux = append(mux, WebController{Post, "POST", "^/"})
	mux = append(mux, WebController{Get, "GET", "^/"})

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", Conf.Port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", Conf.Port))
	log.Fatal(server.ListenAndServe())
}

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	for _, webController := range mux {

		if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m {

			if r.Method == webController.Method {

				webController.Function(w, r)

				go WriteLog(r, t, "match", webController.Pattern)

				return
			}
		}
	}

	go WriteLog(r, t, "unmatch", "")

	io.WriteString(w, "")
	return
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

func Get(w http.ResponseWriter, r *http.Request) {

	client, err := service.NewClient(r, w, Conf.Token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	if len(client.Query.Echostr) > 0 {
		w.Write([]byte(client.Query.Echostr))
		return
	}

	w.WriteHeader(403)
	return
}

func Post(w http.ResponseWriter, r *http.Request) {

	client, err := service.NewClient(r, w, Conf.Token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	client.Run()
	return
}
