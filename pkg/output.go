package pkg

import (
	"bufio"
	"fmt"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"os"
	"sort"
	"strings"
)

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
