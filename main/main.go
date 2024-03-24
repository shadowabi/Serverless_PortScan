package main

import (
	"errors"
	"github.com/shadowabi/Serverless_PortScan_rebuild/cmd"
	"github.com/shadowabi/Serverless_PortScan_rebuild/config"
	"github.com/shadowabi/Serverless_PortScan_rebuild/pkg"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/File"
	"strings"
)

func init() {
	configFile := pkg.GetPwd()
	configFile = strings.Join([]string{configFile, "/config.json"}, "")
	err := File.FileNonExistCreate(configFile)
	Error.HandleFatal(err)
	config.SpecificInit(configFile)
	if config.C.ServerUrl == "http://" || config.C.ServerUrl == "" || config.C.PortList == "" {
		Error.HandleFatal(errors.New("请配置config.json"))
		return
	}
}

func main() {
	cmd.Execute()
}
