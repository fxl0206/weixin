package conf

import (
	"os"
)

const (
	defaultConfFile   string = "conf/yx.ini"
	defaultServerName string = "yxserver"
)

var CurConfFile = defaultConfFile
var CurServerName = defaultServerName

func Initail() {
	args := os.Args
	if len(args) > 1 {
		CurConfFile = args[1]
	}
	if len(args) > 2 {
		CurServerName = args[2]
	}
}
