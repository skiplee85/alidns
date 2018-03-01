package main

import (
	"os"

	"github.com/skiplee85/alidns/base"
)

func main() {

	args := os.Args
	if len(args) < 3 {
		panic("Args < 3, need id,secret,domain.")
	}
	aliDNS := &base.AliDNS{
		AccessKeyID:     args[0],
		AccessKeySecret: args[1],
	}
	ip := base.GetIP()
	rs := aliDNS.GetDomainRecords(args[2], "*")
	if len(rs) > 0 && ip != "" && rs[0].Value != ip {
		rs[0].Value = ip
		aliDNS.UpdateDomainRecord(rs[0])
	}

}
