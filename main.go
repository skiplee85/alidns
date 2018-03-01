package main

import (
	"os"

	"github.com/skiplee85/alidns/base"
)

func main() {

	args := os.Args
	if len(args) < 5 {
		panic("Args < 4, need id,secret,domain,key.")
	}
	aliDNS := &base.AliDNS{
		AccessKeyID:     args[1],
		AccessKeySecret: args[2],
	}
	ip := base.GetIP()
	rs := aliDNS.GetDomainRecords(args[3], args[4])
	if len(rs) > 0 && ip != "" && rs[0].Value != ip {
		rs[0].Value = ip
		aliDNS.UpdateDomainRecord(rs[0])
	}

}
