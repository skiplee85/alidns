package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	baseURL = "http://alidns.aliyuncs.com/"
)

type domainRecordsResp struct {
	RequestID     string `json:"RequestId"`
	TotalCount    int
	PageNumber    int
	PageSize      int
	DomainRecords domainRecords
}

type domainRecords struct {
	Record []DomainRecord
}

// DomainRecord 域名解析记录
type DomainRecord struct {
	DomainName string
	RecordID   string `json:"RecordId"`
	RR         string
	Type       string
	Value      string
	Line       string
	Priority   int
	TTL        int
	Status     string
	Locked     bool
}

func getHTTPBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return body, err
	}
	return nil, fmt.Errorf("Status %d, Error:%s", resp.StatusCode, body)
}

// GetIP 获取公网IP
func GetIP() string {
	var body []byte
	var err error
	// taobao
	body, _ = getHTTPBody("http://www.taobao.com/help/getip.php")
	if err != nil {
		fmt.Printf("GetIP from taobao error.%+v\n", err)
	} else {
		reg := regexp.MustCompile(`"([\d.]+)"`)
		ret := reg.FindStringSubmatch(string(body))
		return ret[1]
	}

	// 360
	body, _ = getHTTPBody("http://ip.360.cn/IPShare/info")
	if err != nil {
		fmt.Printf("GetIP from 360 error.%+v\n", err)
	} else {
		reg := regexp.MustCompile(`"ip":"([\d.]+)"`)
		ret := reg.FindStringSubmatch(string(body))
		return ret[1]
	}

	return ""
}
