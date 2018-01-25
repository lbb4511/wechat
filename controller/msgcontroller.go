package controller

import (
	"log"
	"net/http"

	"github.com/lbb4511/wechat/service"
	"github.com/lbb4511/wechat/setting"
)

func Get(w http.ResponseWriter, r *http.Request) {

	client, err := service.NewClient(r, w, setting.Conf.Token)

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

	client, err := service.NewClient(r, w, setting.Conf.Token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	client.Run()
	return
}
