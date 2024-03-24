package pkg

import (
	"bufio"
	"fmt"
	"github.com/shadowabi/Serverless_PortScan_rebuild/config"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"os"
	"regexp"
	"sort"
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

func GetResult(respData ...ResponseData) (writeResultList []string) {
	if len(respData) != 0 {
		for _, data := range respData {
			dataBody := strings.Trim(data.body, "\"")
			ports := strings.Split(dataBody, ",")
			for _, port := range ports {
				fmt.Println(fmt.Sprintf("[+] %s:%s TCP OPEN.", data.host, port))
				writeResult := strings.Join([]string{data.host, ":", port}, "")
				writeResultList = append(writeResultList, writeResult)
			}
		}
	}
	return writeResultList
}

func WriteToFile(writeResultList []string, output string) {
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	Error.HandlePanic(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	sort.Strings(writeResultList)
	if len(writeResultList) != 0 {
		for _, i := range writeResultList {
			fmt.Fprintln(writer, i)
		}
	}
}
