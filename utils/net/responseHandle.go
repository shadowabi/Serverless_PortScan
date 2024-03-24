package net2

import (
	"bytes"
	"github.com/shadowabi/Serverless_PortScan_rebuild/utils/Error"
	"io"
	"net/http"
)

func HandleResponse(resp *http.Response) (bodyString string) {
	bodyBuf := new(bytes.Buffer)
	_, err := io.Copy(bodyBuf, resp.Body)
	Error.HandleError(err)
	bodyString = bodyBuf.String()
	return bodyString
}
