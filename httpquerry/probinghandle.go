package httpquerry

import (
	"net/http"
	"redissyncer-portal/global"
	"time"
)

//探活，通过向节点health接口发送请求判断节点是否存活
func NodeAlive(addr, port string) bool {

	httpclient := &http.Client{}
	httpclient.Timeout = 10 * time.Second
	url := "http://" + addr + ":" + port + "/health"

	global.RSPLog.Sugar().Debug(url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return false
	}

	resp, respErr := httpclient.Do(req)
	if respErr != nil {
		global.RSPLog.Sugar().Debug(respErr)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		global.RSPLog.Sugar().Debug(resp.StatusCode)
		return false
	}

	return true
}
