package cmd

import (
	"errors"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/shadowabi/Serverless_PortScan_rebuild/config"
	"github.com/shadowabi/Serverless_PortScan_rebuild/define"
	"github.com/shadowabi/Serverless_PortScan_rebuild/pkg"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/log"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "Serverless_PortScan",
	Short: "Serverless_PortScan is used to scan ports using cloud functions.",
	Long: "  ____                           _                   ____            _   ____                  \n" +
		" / ___|  ___ _ ____   _____ _ __| | ___  ___ ___    |  _ \\ ___  _ __| |_/ ___|  ___ __ _ _ __  \n" +
		" \\___ \\ / _ \\ '__\\ \\ / / _ \\ '__| |/ _ \\/ __/ __|   | |_) / _ \\| '__| __\\___ \\ / __/ _` | '_ \\ \n" +
		"  ___) |  __/ |   \\ V /  __/ |  | |  __/\\__ \\__ \\   |  __/ (_) | |  | |_ ___) | (_| (_| | | | | \n" +
		" |____/ \\___|_|    \\_/ \\___|_|  |_|\\___||___/___/___|_|   \\___/|_|   \\__|____/ \\___\\__,_|_| |_| " +
		"                                               	  			|_____|                                          \n" +
		` 
        github.com/shadowabi/Serverless_PortScan

Serverless_PortScan是一个云函数端口扫描器。
Serverless_PortScan is used to scan ports using cloud functions.
`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Init(logLevel)
		if define.Url != "" && define.File != "" {
			Error.HandleFatal(errors.New("参数不可以同时存在"))
			return
		}
		if define.Url == "" && define.File == "" {
			Error.HandleFatal(errors.New("必选参数为空，请输入 -u 参数或 -f 参数"))
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if define.Port == "" {
			define.Port = config.C.PortList
		}
		var hostList []string
		if define.File != "" {
			hostList = pkg.ParseFileParameter(define.File)
		} else {
			hostList = append(hostList, define.Url)
		}

		hostList = pkg.RemoveDuplicates(hostList)
		client := pkg.GenerateHTTPClient(define.TimeOut * len(hostList)) // 动态分配 http 超时时间

		reqList := pkg.ConvertToReqList(define.Port, hostList...)
		resp := pkg.FetchPortData(client, reqList...)

		writeList := pkg.GetResult(resp...)
		pkg.WriteToFile(writeList, define.OutPUT)
	},
}

var logLevel string

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "设置日志等级 (Set log level) [trace|debug|info|warn|error|fatal|panic]")
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().StringVarP(&define.File, "file", "f", "", "从文件中读取目标地址 (Input filename)")
	RootCmd.Flags().StringVarP(&define.Url, "url", "u", "", "输入目标地址 (Input [ip|domain|url])")
	RootCmd.Flags().StringVarP(&define.Port, "port", "p", "", "输入需要被扫描的端口，逗号分割 (Enter the port to be scanned, separated by commas (,))")
	RootCmd.Flags().IntVarP(&define.TimeOut, "timeout", "t", 5, "输入每个 http 请求的超时时间 (Enter the timeout period for every http request)")
	RootCmd.Flags().StringVarP(&define.OutPUT, "output", "o", "./result.txt", "输入结果文件输出的位置 (Enter the location of the scan result output)")
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiGreen + cc.Underline,
		Commands: cc.Cyan + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Cyan + cc.Bold,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
