package main

import (
	"conf"
	"yxfs.com/yxserver"
)

func main() {
	conf.Initail()
	server := yxserver.NewInstance(conf.CurConfFile, conf.CurServerName)
	server.Start()
}
