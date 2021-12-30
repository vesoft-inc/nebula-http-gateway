package main

import (
	"log"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/gateway/dao"
)

func main() {
	var (
		address  = "192.168.8.169"
		port     = 9669
		username = "root"
		password = "123"
	)

	nsid, err := dao.Connect(address, port, username, password, nebula.VersionAuto)
	if err != nil {
		log.Println("error: ", err)
	}
	log.Println(nsid)
	defer dao.Disconnect(nsid)

	gql := "show hosts;"
	res, p, err := dao.Execute(nsid, gql)
	if err != nil {
		log.Println("error: ", err)
		if p != nil {
			log.Fatal(p)
		}
	}

	log.Println(res)
}
