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

	nsid, err := dao.Connect(address, port, username, password, nebula.VersionAuto)
	if err != nil {
		log.Println("error: ", err)
	}
	log.Println(nsid)
	defer dao.Disconnect(nsid)

	gql := "CREATE SPACE IF NOT EXISTS basic_example_space(vid_type=FIXED_STRING(20)); " +
		"USE basic_example_space;" +
		"CREATE TAG IF NOT EXISTS person(name string, age int);" +
		"CREATE EDGE IF NOT EXISTS like(likeness double);"
	res, p, err := dao.Execute(nsid, gql, nil)
	if err != nil {
		log.Println("error: ", err)
		if p != nil {
			log.Fatal(p)
		}
	}

	log.Println(res)

	gql = "DROP SPACE IF EXISTS basic_example_space;"
	res, p, err = dao.Execute(nsid, gql, nil)
	if err != nil {
		log.Println("error: ", err)
		if p != nil {
			log.Fatal(p)
		}
	}

	log.Println(res)
}
