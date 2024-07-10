package pkg

import (
	"crypto/tls"
	"github.com/shadowabi/Serverless_PortScan_rebuild/define"
	"net/http"
	"time"
)

func GenerateHTTPClient(timeOut int) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return client
}

func IsEmptyConfig(c define.Configure) bool {
	return c.ServerUrl == "" ||
		c.PortList == ""
}
