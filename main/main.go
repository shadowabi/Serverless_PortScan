package main

import (
	"fmt"

	"github.com/shadowabi/Serverless_PortScan_rebuild/cmd"
	"github.com/shadowabi/Serverless_PortScan_rebuild/config"
	"github.com/shadowabi/Serverless_PortScan_rebuild/pkg"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
)

func init() {
	config.InitConfigure("config.yaml")
}

func main() {
	if pkg.IsEmptyConfig(config.C) {
		Error.HandleFatal(fmt.Errorf("请配置 config.yaml"))
		return
	}

	fmt.Println(cmd.RootCmd.Long)
	cmd.Execute()
}
