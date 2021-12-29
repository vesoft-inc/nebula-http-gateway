package main

import (
	"log"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/gateway/dao"
)

func main() {
	var (
		address  = "192.168.8.167"
		port     = 9669
		username = "root"
		password = "123"
	)
	version, err := nebula.VersionHelper(address, port, username, password)
	if err != nil {
		log.Println("[Error]", err)
		version = "v2.5"
	}
	log.Println(version)

	nsid, err := dao.Connect(address, port, username, password, string(version))
	if err != nil {
		log.Println("error: ", err)
	}
	log.Println(nsid)
	defer dao.Disconnect(nsid)

	gql := "show users;"
	res, p, err := dao.Execute(nsid, gql)
	if err != nil {
		log.Println("error: ", err)
		if p != nil {
			log.Fatal(p)
		}
	}

	log.Println(res)
}
