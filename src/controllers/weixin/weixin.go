package weixin

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/sidbusy/weixinmp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

//weixin Controller 请求路径
//var Prefix string = "/sign"
var SignPrefix string = ""
var ActionPrefix string = "/action"

//微信请求消息结构
type Request struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

//微信响应头结构
type Response struct {
	ToUserName   string `xml:"xml>ToUserName"`
	FromUserName string `xml:"xml>FromUserName"`
	CreateTime   string `xml:"xml>CreateTime"`
	MsgType      string `xml:"xml>MsgType"`
	Content      string `xml:"xml>Content"`
	MsgId        int    `xml:"xml>MsgId"`
}

//服务器Token认证
func DoSign(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(SignPrefix):]
	log.Println(path)
	r.ParseForm()
	token := "myluckyfxl"
	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	var echostr string = strings.Join(r.Form["echostr"], "")
	fmt.Println("signature:" + signature)
	fmt.Println("timestamp:" + timestamp)
	fmt.Println("nonce:" + nonce)
	fmt.Println("echostr:" + echostr)
	tmps := []string{token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	fmt.Println("tmp:" + tmp)
	if tmp == signature {
		log.Println("signature Success")
		fmt.Fprintf(w, echostr)
	} else {
		log.Println("signature Failed!")
		fmt.Fprintf(w, echostr)
	}
}

//消息处理
func DoAction(w http.ResponseWriter, r *http.Request) {
	log.Println("here")

	postedMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()
	v := Request{}
	xml.Unmarshal(postedMsg, &v)
	if v.MsgType == "text" {
		v := Request{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, v.Content, v.MsgId}
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Fprintf(w, string(output))
	} else if v.MsgType == "event" {
		Content := `"欢迎关注  我的微信"`
		v := Request{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, Content, v.MsgId}
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Fprintf(w, string(output))
	} else {
		log.Println("Something is Wrong!")
	}
}
func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//微信方式2
func Receiver(w http.ResponseWriter, r *http.Request) {
	token := "myluckyfxl"                        // 微信公众平台的Token
	appid := "wxe23370ffc2050dba"                // 微信公众平台的AppID
	secret := "99a335cf1a73c15cb50422dd3004fd41" // 微信公众平台的AppSecret
	// 仅被动响应消息时可不填写appid、secret
	// 仅主动发送消息时可不填写token
	mp := weixinmp.New(token, appid, secret)
	// 检查请求是否有效
	// 仅主动发送消息时不用检查
	if !mp.Request.IsValid(w, r) {
		return
	}
	// 判断消息类型
	if mp.Request.MsgType == weixinmp.MsgTypeText {
		fmt.Println("message :!" + mp.Request.Content)
		// 回复消息
		mp.ReplyTextMsg(w, "自动回复：lucky 好哈都不晓得！ 你说了 "+mp.Request.Content)
		var err=mp.SendTextMsg("fxl0206", mp.Request.Content)\
		if err!=nil {
			fmt.Print(err)
		}
	} else if mp.Request.MsgType == weixinmp.MsgTypeEvent {
		mp.ReplyTextMsg(w, "感谢支持关注灰色的小鸟！")
	}
}
