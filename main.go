package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lbb4511/wechat/wx"
)

const (
	logLevel = "dev"
	port     = 8098
	token    = "asgdyg3wyw4qfasd2gyw3"
)

func get(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

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

func post(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	client.Run()
	return
}

func main() {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", port))
	log.Fatal(server.ListenAndServe())
}
