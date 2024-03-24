package pkg

import (
	"crypto/tls"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"net/http"
	"os"
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

func GetPwd() (homePath string) {
	homePath, err := os.Getwd()
	Error.HandlePanic(err)
	return homePath
}
