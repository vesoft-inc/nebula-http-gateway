package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/wrapper"
)

func main() {
	for _, version := range []nebula.Version{
		nebula.Version2_5,
		nebula.Version2_6,
		nebula.Version3_0,
	} {
		var (
			c   nebula.Client
			gc  nebula.GraphClient
			err error

			host     = "192.168.8.167:9669"
			username = "root"
			password = "123"
		)
		{ // use nebula.NewClient
			c, err = nebula.NewClient(nebula.ConnectionInfo{
				GraphEndpoints: []string{host},
				GraphAccount: nebula.Account{
					Username: username,
					Password: password,
				},
			}, nebula.WithVersion(version))
			if err != nil {
				panic(fmt.Sprintf("%s %+v", version, err))
			}
			if err = c.Graph().Open(); err != nil {
				panic(err)
			}

			executeExample(c.Graph())

			if err = c.Graph().Close(); err != nil {
				panic(err)
			}
			log.Println("basic execute example finished")
		}
		{ // use nebula.NewGraphClient
			log.Println("execute with params example...")
			gc, err = nebula.NewGraphClient([]string{host}, username, password, nebula.WithVersion(version))
			if err != nil {
				panic(fmt.Sprintf("%s %+v", version, err))
			}
			if err = gc.Open(); err != nil {
				panic(err)
			}

			executeWithParamsExample(gc)

			if err = gc.Close(); err != nil {
				panic(err)
			}
			log.Println("execute with params example finished")
		}

		factoryExample(c)
	}
}

func executeExample(gc nebula.GraphClient) {
	{
		// Prepare the query
		createSchema := "CREATE SPACE IF NOT EXISTS basic_example_space(vid_type=FIXED_STRING(20)); " +
			"USE basic_example_space;" +
			"CREATE TAG IF NOT EXISTS person(name string, age int);" +
			"CREATE EDGE IF NOT EXISTS like(likeness double);"

		// Excute a query
		resp, err := gc.Execute([]byte(createSchema))
		if err != nil {
			panic(err)
		}

		resultSet, _ := wrapper.GenResultSet(resp, gc.Factory(), types.TimezoneInfo{})
		checkResultSet(createSchema, resultSet)
		log.Println(resp)
	}
	// Drop space
	{
		query := "DROP SPACE IF EXISTS basic_example_space;"
		// Send query
		resp, err := gc.Execute([]byte(query))
		if err != nil {
			panic(err)
		}

		resultSet, _ := wrapper.GenResultSet(resp, gc.Factory(), types.TimezoneInfo{})
		checkResultSet(query, resultSet)
	}
}

func executeWithParamsExample(gc nebula.GraphClient) {
	params := make(map[string]interface{})
	params["p1"] = true
	params["p2"] = 3
	params["p3"] = []interface{}{true, 3}
	params["p4"] = map[string]interface{}{"a": true, "b": 3}

	gql := "RETURN abs($p2)+1 AS col1, toBoolean($p1) and false AS col2, $p3, $p4.a"
	resp, err := gc.ExecuteWithParameter([]byte(gql), params)
	if err != nil {
		log.Printf("execute with params failed: %s\n", err.Error())
		return
	}

	resultSet, _ := wrapper.GenResultSet(resp, gc.Factory(), types.TimezoneInfo{})
	checkResultSet(gql, resultSet)
	// Get all column names from the resultSet
	colNames := resultSet.GetColNames()
	fmt.Printf("Column names: %s\n", strings.Join(colNames, ", "))
	fmt.Println(resultSet.AsStringTable())
	// Get a row from resultSet
	record, err := resultSet.GetRowValuesByIndex(0)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Print whole row
	fmt.Printf("The first row elements: %s\n", record.String())
}

func factoryExample(c nebula.Client) {
	examplef := func(factory nebula.Factory) {
		s := []byte{1, 2, 3}
		vb := factory.NewValueBuilder()
		vb.SVal(s)
		v1 := vb.Build()
		v2 := vb.Build().SetSVal([]byte{1, 2})
		i1 := v1.Unwrap()
		i2 := v2.Unwrap()
		log.Printf("\n%v, %v\n%p, %p;\n%v, %v", v1 == v2, i1 == i2, i1, i2, v1.GetSVal(), v2.GetSVal())
	}

	log.Println("get factory by client version:")
	factory, _ := nebula.NewFactory(nebula.WithVersion(c.Version()))
	examplef(factory)

	log.Println("get factory by client:")
	examplef(c.Factory())

	log.Println("get factory by graph client:")
	examplef(c.Graph().Factory())
}

func checkResultSet(prefix string, res *wrapper.ResultSet) {
	if !res.IsSucceed() {
		log.Fatal(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
	}
}
