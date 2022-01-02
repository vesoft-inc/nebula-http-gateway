package main

import (
	"fmt"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"log"
)

func main() {
	for _, version := range []nebula.Version{
		nebula.Version2_5,
		nebula.Version2_6,
		nebula.Version3_0,
	} {
		{ // use nebula.NewClient
			c, err := nebula.NewClient(nebula.ConnectionInfo{
				GraphEndpoints: []string{"192.168.8.167:9669"},
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
		{ // test factory and builder
			f, err := nebula.NewFactory(nebula.WithVersion(version))
			if err != nil {
				log.Fatal(err)
			}

			s := []byte{1, 2, 3}
			v1 := f.NewValue()
			v1.SetSVal(s)
			v2 := v1.Builder().Emit().SetSVal([]byte{1, 2})
			i1 := v1.Unwrap()
			i2 := v2.Unwrap()
			log.Printf("\n%p, %p;\n%v, %v", i1, i2, v1.GetSVal(), v2.GetSVal())
		}
	}

}
