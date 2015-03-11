package yxserver

import (
	"controllers/weixin"
	"github.com/widuu"
	"log"
	"net/http"
)

//全局配置文路径
var ConfPath string

//websocket请求路径
var socketPath string

//系统配置初始化
func NewInstance(confPath string, serverName string) WbServer {
	conf := goini.SetConfig(confPath)
	port := conf.GetValue(serverName, "port")
	viewPath := conf.GetValue(serverName, "viewPath")
	ConfPath = confPath
	wsPath := conf.GetValue("websocket", "wsPath")
	wsPort := conf.GetValue("websocket", "port")
	wsIp := conf.GetValue("websocket", "wsIp")
	socketPath = "ws://" + wsIp + ":" + wsPort + "/" + wsPath
	return WbServer{port, viewPath, serverName}
}

type WbServer struct {
	port       string
	viewPath   string
	ServerName string
}

//启动服务器，绑定请求控制器
func (this *WbServer) Start() {
	//微信token验证
	http.HandleFunc(weixin.SignPrefix+"/", weixin.DoSign)
	//微信消息处理
	//http.HandleFunc(weixin.ActionPrefix+"/", weixin.DoAction)
	http.HandleFunc("/wx", weixin.Receiver)
	//监听端口
	err := http.ListenAndServe(":"+this.port, nil)
	//服务器启动错误处理
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
