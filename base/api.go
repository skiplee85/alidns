package base

// https://help.aliyun.com/document_detail/29739.html
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// AliDNS 阿里云token
type AliDNS struct {
	AccessKeyID     string
	AccessKeySecret string
}

var publicParm = map[string]string{
	"AccessKeyId":      "",
	"Format":           "JSON",
	"Version":          "2015-01-09",
	"SignatureMethod":  "HMAC-SHA1",
	"Timestamp":        "",
	"SignatureVersion": "1.0",
	"SignatureNonce":   "",
}

// GetDomainRecords 获取解析记录列表
func (d *AliDNS) GetDomainRecords(domain, rr string) []DomainRecord {
	resp := &domainRecordsResp{}
	parms := map[string]string{
		"Action":     "DescribeDomainRecords",
		"DomainName": domain,
		"RRKeyWord":  rr,
	}
	urlPath := d.genRequestURL(parms)
	body, err := getHTTPBody(urlPath)
	if err != nil {
		fmt.Printf("GetDomainRecords error.%+v\n", err)
	} else {
		json.Unmarshal(body, resp)
		return resp.DomainRecords.Record
	}
	return nil
}

// UpdateDomainRecord 修改解析记录
func (d *AliDNS) UpdateDomainRecord(r DomainRecord) error {
	parms := map[string]string{
		"Action":   "UpdateDomainRecord",
		"RecordId": r.RecordID,
		"RR":       r.RR,
		"Type":     r.Type,
		"Value":    r.Value,
		"TTL":      strconv.Itoa(r.TTL),
		"Line":     r.Line,
	}

	urlPath := d.genRequestURL(parms)
	body, err := getHTTPBody(urlPath)
	if err != nil {
		fmt.Printf("UpdateDomainRecord error.%+v\n", err)
	} else {
		fmt.Printf("UpdateDomainRecord succ. %s\n", body)
	}
	return err
}

func (d *AliDNS) genRequestURL(parms map[string]string) string {
	pArr := []string{}
	ps := map[string]string{}
	for k, v := range publicParm {
		ps[k] = v
	}
	for k, v := range parms {
		ps[k] = v
	}
	now := time.Now().UTC()
	ps["AccessKeyId"] = d.AccessKeyID
	ps["SignatureNonce"] = strconv.Itoa(int(now.UnixNano()) + rand.Intn(99999))
	ps["Timestamp"] = now.Format("2006-01-02T15:04:05Z")

	for k, v := range ps {
		pArr = append(pArr, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(pArr)
	path := strings.Join(pArr, "&")

	s := "GET&%2F&" + url.QueryEscape(path)
	s = strings.Replace(s, "%3A", "%253A", -1)
	s = strings.Replace(s, "%40", "%2540", -1)
	s = strings.Replace(s, "%2A", "%252A", -1)
	mac := hmac.New(sha1.New, []byte(d.AccessKeySecret+"&"))

	mac.Write([]byte(s))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s?%s&Signature=%s", baseURL, path, url.QueryEscape(sign))
}
