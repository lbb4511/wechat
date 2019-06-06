package service

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/clbanning/mxj"
	"github.com/lbb4511/wechat/models"
)

// client
type client models.WeixinClient

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

func value2CDATA(v string) CDATAText {
	return CDATAText{"<![CDATA[" + v + "]]>"}
}

func (b *Base) InitBaseData(w *client, msgtype string) {

	b.FromUserName = value2CDATA(w.Message["ToUserName"].(string))
	b.ToUserName = value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = value2CDATA(strconv.FormatInt(time.Now().Unix(), 10))
	b.MsgType = value2CDATA(msgtype)
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*client, error) {

	cl := new(client)

	cl.Token = token
	cl.Request = r
	cl.ResponseWriter = w

	cl.initWeixinQuery()

	if cl.Query.Signature != cl.signature() {
		return nil, errors.New("Invalid Signature.")
	}

	return cl, nil
}

func (this *client) initWeixinQuery() {

	var q models.WeixinQuery

	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")

	this.Query = q
}

func (this *client) signature() string {

	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *client) initMessage() error {

	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

func (this *client) text() {

	inMsg, ok := this.Message["Content"].(string)

	if !ok {
		return
	}

	var reply TextMessage

	reply.InitBaseData(this, "text")
	reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))

	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

func (this *client) Run() {

	err := this.initMessage()

	if err != nil {

		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		this.text()
		break
	default:
		break
	}

	return
}
