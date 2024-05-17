package pkg

import (
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	net2 "github.com/shadowabi/Serverless_PortScan_rebuild/utils/net"
	"net/http"
	"sync"
	"time"
)

type ResponseData struct {
	host string
	body string
}

func FetchPortData(client *http.Client, reqList ...ReqList) (respData []ResponseData) {
	if len(reqList) != 0 {
		var wg sync.WaitGroup
		respChan := make(chan ResponseData, cap(reqList))
		for _, request := range reqList {
			wg.Add(1)
			time.Sleep(100 * time.Millisecond)
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
		for resp := range respChan {
			respData = append(respData, resp)
		}
	}
	return respData
}
