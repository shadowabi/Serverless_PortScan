package pkg

import (
	"bufio"
	"github.com/shadowabi/Serverless_PortScan_rebuild/config"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"os"
	"regexp"
	"strings"
)

type ReqList struct {
	req  string
	host string
}

func ParseFileParameter(fileName string) (fileHostList []string) {
	file, err := os.Open(fileName)
	Error.HandlePanic(err)
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		fileHostList = append(fileHostList, line)
	}
	file.Close()
	return fileHostList
}

func ConvertToReqList(port string, param ...string) (reqList []ReqList) {
	if len(param) != 0 {
		for _, host := range param {
			// 当输入Url时提取出域名
			re := regexp.MustCompile(`(http|https)://`)
			if re.MatchString(host) {
				host = re.ReplaceAllString(host, "")
				re = regexp.MustCompile(`([/\\]).*`)
				if re.MatchString(host) {
					host = re.ReplaceAllString(host, "")
				}
			}
			// 匹配IP/域名
			if regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`).MatchString(host) ||
				regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$`).MatchString(host) {
				target := strings.Join([]string{config.C.ServerUrl, "?ip=", host, "&port=", port}, "")
				reqList = append(reqList, ReqList{target, host})
			}
		}
	}
	return reqList
}

func RemoveDuplicates(reqList []string) []string {
	encountered := map[string]bool{}
	var result []string
	for v := range reqList {
		if !encountered[reqList[v]] {
			encountered[reqList[v]] = true
			result = append(result, reqList[v])
		}
	}
	return result
}
