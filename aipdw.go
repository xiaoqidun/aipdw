package aipdw

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	client  = &http.Client{Timeout: 5 * time.Second}
	ipCache = &sync.Map{}
)

type Body struct {
	Status   string `json:"status"`   // 返回结果状态值，值为0或1,0表示失败；1表示成功
	Info     string `json:"info"`     // status为0时，info返回错误原因；否则返回“OK”
	Infocode string `json:"infocode"` // 返回状态说明,10000代表正确,详情参阅info状态表
	Result
}

type Result struct {
	Country  string `json:"country"`  // 国家
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	District string `json:"district"` // 区县
	Isp      string `json:"isp"`      // 运营商
	Location string `json:"location"` // 经纬度
	IP       string `json:"ip"`       // IP地址
}

// QueryIP 使用高德位置服务查询IP
func QueryIP(sk string, key string, ip string) (result Result, err error) {
	v, ok := ipCache.Load(ip)
	if ok {
		result = v.(Result)
		return
	}
	pip := net.ParseIP(ip)
	if pip == nil {
		err = errors.New("ip is not ipv4 or ipv6")
		return
	}
	ipt := "6"
	if pip.To4() != nil {
		ipt = "4"
	}
	arg := &reqLBS{
		SK: sk,
		Args: map[string]string{
			"ip":   ip,
			"key":  key,
			"type": ipt,
		},
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", "https://restapi.amap.com/v5/ip", arg.Encode()), nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var bodyUnmarshal Body
	err = json.Unmarshal(body, &bodyUnmarshal)
	if err != nil {
		return
	}
	if bodyUnmarshal.Infocode != "10000" {
		err = fmt.Errorf("resp code is %s, body is %s", bodyUnmarshal.Infocode, body)
		return
	}
	result = bodyUnmarshal.Result
	ipCache.Store(ip, bodyUnmarshal.Result)
	return
}
