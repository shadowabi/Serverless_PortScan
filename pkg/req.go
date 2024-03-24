package pkg

import (
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	net2 "github.com/shadowabi/Serverless_PortScan_rebuild/utils/net"
	"net/http"
	"sync"
)

type ResponseData struct {
	host string
	body string
}

func FetchPortData(client *http.Client, reqList ...ReqList) (respData []ResponseData) {
	if len(reqList) != 0 {
		var wg sync.WaitGroup
		respChan := make(chan ResponseData, len(reqList))
		for _, request := range reqList {
			wg.Add(1)
			go func(request ReqList, wg *sync.WaitGroup) {
				defer wg.Done()
				resp, err := client.Get(request.req)
				Error.HandleError(err)
				respBody := net2.HandleResponse(resp)
				resp.Body.Close()
				respChan <- ResponseData{request.host, respBody}
			}(request, &wg)
		}
		wg.Wait()
		close(respChan)

		for data := range respChan { // 从channel中收集结果
			respData = append(respData, data)
		}
	}
	return respData
}
