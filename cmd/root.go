package cmd

import (
	"os"
	"github.com/spf13/cobra"
	exec "github.com/shadowabi/Serverless_PortScan/pkg"
	"strings"
	"sync"
)

func Execute(){
	var (
		file	string
		url		string
		port 	string
	)
	exec.ReadConfig()

	var RootCmd = &cobra.Command{
		Use:   "Serverless_PortScan",
		Short: "Serverless_PortScan is used to scan ports using cloud functions.",
		Long: 
	"  ____                           _                   ____            _   ____                  \n" +
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
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
            if url != "" {
            	wg.Add(1)
                go exec.CheckIP(strings.TrimSpace(url), &wg)
                wg.Wait()
            } else if file != "" {
                exec.ReadFile(file)
            }

            if len(exec.Url) != 0 {
            	exec.Scan(port)
            	exec.OutPut()
            }
        },
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().StringVarP(&file, "file", "f", "", "从文件中读取目标地址 (Input FILENAME)")
	RootCmd.Flags().StringVarP(&url, "url", "u", "", "输入目标地址 (Input IP/DOMAIN/URL)")
	RootCmd.Flags().StringVarP(&port, "port", "p", exec.Config.Default_port , "从文件中读取网站地址 (Input FILENAME)")

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}