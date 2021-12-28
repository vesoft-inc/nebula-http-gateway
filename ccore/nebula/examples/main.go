package main

import (
	"fmt"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
)

func main() {
	for _, version := range []nebula.Version{
		nebula.V2_5,
		nebula.V2_6,
		nebula.V3_0,
	} {
		{ // use nebula.NewClient
			c, err := nebula.NewClient(nebula.ConnectionInfo{
				GraphEndpoints: []string{"192.168.8.169:9669"},
				GraphAccount: nebula.Account{
					Username: "root",
					Password: "123",
				},
			}, nebula.WithVersion(version))
			if err != nil {
				panic(fmt.Sprintf("%s %+v", version, err))
			}
			if err := c.Graph().Open(); err != nil {
				panic(err)
			}
			resp, err := c.Graph().Execute([]byte("show users;"))
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", resp)
			if err := c.Graph().Close(); err != nil {
				panic(err)
			}
		}
		{ // use nebula.NewClient
			c, err := nebula.NewGraphClient([]string{"192.168.8.169:9669"}, "root", "123", nebula.WithVersion(version))
			if err != nil {
				panic(fmt.Sprintf("%s %+v", version, err))
			}
			if err := c.Open(); err != nil {
				panic(err)
			}
			resp, err := c.Execute([]byte("show users;"))
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", resp)
			if err := c.Close(); err != nil {
				panic(err)
			}
		}
	}
}
