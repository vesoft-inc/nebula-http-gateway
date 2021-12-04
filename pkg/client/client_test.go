package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_0_0_ga"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_0"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_1"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_6_0"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/logger"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/types"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var (
	address = types.HostAddress{
		Host: "192.168.8.49",
		Port: 9669,
	}
	account = types.Account{
		Username: "root",
		Password: "123456",
	}
	config = types.ClientConfig{
		TimeOut:         0 * time.Millisecond,
		IdleTime:        0 * time.Millisecond,
		MaxConnPoolSize: 10,
		MinConnPoolSize: 0,
		SslConfig:       nil,
	}

	rootCA     = openAndReadFile("./tmp/test.ca.pem")
	cert       = openAndReadFile("./tmp/test.client.crt")
	privateKey = openAndReadFile("./tmp/test.client.key")

	sslConfig           = testSslConfig()
	configConfigWithSSL = types.ClientConfig{
		TimeOut:         0 * time.Millisecond,
		IdleTime:        0 * time.Millisecond,
		MaxConnPoolSize: 10,
		MinConnPoolSize: 0,
		SslConfig:       sslConfig,
	}
	stmts = []string{
		"drop space test;",
	}
)

func TestV2_0_0_ga(t *testing.T) {
	log.Println("\n=> testing v2.0.0-ga")
	config.Ver = types.NewVersion("v2.0.0-ga")
	fmt.Printf("test service address: %s:%d\n", address.Host, address.Port)
	c, _ := NewClient([]types.HostAddress{address}, config, logger.DefaultLogger{})
	log.Println("version test pass", c.Version())

	err := c.Ping(address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
	log.Println("ping test pass")

	s, err := c.NewSession(account)
	if err != nil {
		t.Error(err)
	}
	log.Println("new session test pass")

	for _, stmt := range stmts {
		fmt.Println("exec test stmt: ", stmt)
		rset, err := s.Execute(stmt)
		if err != nil {
			t.Error(err)
		}
		resultSet, err := clientv2_0_0_ga.ResultSetWrapper(rset)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("is result set succeed: ", resultSet.IsSucceed())
	}
	log.Println("execute test pass")

	_, err = s.ExecuteJson(stmts[0])
	if err != nil {
		log.Println(err)
		log.Println("execute json test pass")
	} else {
		t.Error("execute json test failed")
	}

	s.Release() // test client sessionPool status
	_, err = s.Execute(stmts[0])
	if err != nil {
		log.Println("release session test pass")
	} else {
		t.Error("session release test failed")
	}

	c.Close()
	log.Println("client close test pass")
}

func TestV2_5_0(t *testing.T) {
	log.Println("\n=> testing v2.5.0")
	config.Ver = types.NewVersion("v2.5.0")
	fmt.Printf("test service address: %s:%d\n", address.Host, address.Port)
	c, _ := NewClient([]types.HostAddress{address}, config, logger.DefaultLogger{})
	log.Println("version test pass", c.Version())

	err := c.Ping(address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
	log.Println("ping test pass")

	s, err := c.NewSession(account)
	if err != nil {
		t.Error(err)
	}
	log.Println("new session test pass")

	for _, stmt := range stmts {
		fmt.Println("exec test stmt: ", stmt)
		rset, err := s.Execute(stmt)
		if err != nil {
			t.Error(err)
		}
		resultSet, err := clientv2_5_0.ResultSetWrapper(rset)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("is result set succeed: ", resultSet.IsSucceed())
	}
	log.Println("execute test pass")

	_, err = s.ExecuteJson(stmts[0])
	if err != nil {
		log.Println(err)
		log.Println("execute json test pass")
	} else {
		t.Error("execute json test failed")
	}

	s.Release() // test client sessionPool status
	_, err = s.Execute(stmts[0])
	if err != nil {
		log.Println("release session test pass")
	} else {
		t.Error("session release test failed")
	}

	c.Close()
	log.Println("client close test pass")
}

func TestV2_5_1(t *testing.T) {
	log.Println("\n=> testing v2.5.1")
	config.Ver = types.NewVersion("v2.5.1")
	fmt.Printf("test service address: %s:%d\n", address.Host, address.Port)
	c, _ := NewClient([]types.HostAddress{address}, config, logger.DefaultLogger{})
	log.Println("version test pass", c.Version())

	err := c.Ping(address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
	log.Println("ping test pass")

	s, err := c.NewSession(account)
	if err != nil {
		t.Error(err)
	}
	log.Println("new session test pass")

	for _, stmt := range stmts {
		fmt.Println("exec test stmt: ", stmt)
		rset, err := s.Execute(stmt)
		if err != nil {
			t.Error(err)
		}
		resultSet, err := clientv2_5_1.ResultSetWrapper(rset)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("is result set succeed: ", resultSet.IsSucceed())
	}
	log.Println("execute test pass")

	_, err = s.ExecuteJson(stmts[0])
	if err != nil {
		log.Println(err)
		log.Println("execute json test pass")
	} else {
		t.Error("execute json test failed")
	}

	s.Release() // test client sessionPool status
	_, err = s.Execute(stmts[0])
	if err != nil {
		log.Println("release session test pass")
	} else {
		t.Error("session release test failed")
	}

	c.Close()
	log.Println("client close test pass")
}

func TestV2_6_0(t *testing.T) {
	log.Println("\n=> testing v2.6.0")
	config.Ver = types.NewVersion("v2.6.0")
	fmt.Printf("test service address: %s:%d\n", address.Host, address.Port)
	c, _ := NewClient([]types.HostAddress{address}, config, logger.DefaultLogger{})
	log.Println("version test pass", c.Version())

	err := c.Ping(address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
	log.Println("ping test pass")

	s, err := c.NewSession(account)
	if err != nil {
		t.Error(err)
	}
	log.Println("new session test pass")

	for _, stmt := range stmts {
		fmt.Println("exec test stmt: ", stmt)
		rset, err := s.Execute(stmt)
		if err != nil {
			t.Error(err)
		}
		resultSet, err := clientv2_6_0.ResultSetWrapper(rset)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("is result set succeed: ", resultSet.IsSucceed())
	}
	log.Println("execute test pass")

	resultAsBytes, err := s.ExecuteJson(stmts[0])
	fmt.Println("result set bytes (str): ", string(resultAsBytes))
	if err != nil {
		t.Error(err)
	}
	log.Println("execute json test pass")

	s.Release() // test client sessionPool status
	_, err = s.Execute(stmts[0])
	if err != nil {
		log.Println("release session test pass")
	} else {
		t.Error("session release test failed")
	}

	c.Close()
	log.Println("client close test pass")
}

func TestV2_6_0_SSL(t *testing.T) {
	log.Println("\n=> testing v2.6.0 with ssl")
	configConfigWithSSL.Ver = types.NewVersion("v2.6.0")
	fmt.Printf("test service address: %s:%d\n", address.Host, address.Port)

	fmt.Println(configConfigWithSSL.SslConfig)
	c, _ := NewClient([]types.HostAddress{address}, configConfigWithSSL, logger.DefaultLogger{})
	log.Println("version test pass", c.Version())

	err := c.Ping(address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
	log.Println("ping test pass")

	s, err := c.NewSession(account)
	if err != nil {
		t.Error(err)
	}
	log.Println("new session test pass")

	stmt := "drop space test;"
	fmt.Println("exec test stmt: ", stmt)
	rset, err := s.Execute(stmt)
	if err != nil {
		t.Error(err)
	}
	resultSet, err := clientv2_6_0.ResultSetWrapper(rset)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("is result set succeed: ", resultSet.IsSucceed())
	log.Println("execute test pass")
	resultAsBytes, err := s.ExecuteJson(stmt)
	fmt.Println("result set bytes (str): ", string(resultAsBytes))
	if err != nil {
		t.Error(err)
	}
	log.Println("execute json test pass")

	s.Release() // test client sessionPool status
	rset, err = s.Execute(stmts[0])
	if err != nil {
		log.Println("release session test pass")
	} else {
		t.Error("session release test failed")
	}

	c.Close()
	log.Println("client close test pass")
}

func testSslConfig() *tls.Config {
	// generate the client certificate
	clientCert, err := tls.X509KeyPair(cert, privateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to get key pair%v", err))
	}

	// parse root CA pem and add into CA pool
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(rootCA)
	if !ok {
		panic("unable to append supplied cert into tls.Config, are you sure it is a valid certificate")
	}

	// set tls config
	// InsecureSkipVerify is set to true for test purpose ONLY. DO NOT use it in production.
	sslConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            rootCAPool,
		InsecureSkipVerify: true, // This is only used for testing
	}

	return sslConfig
}

func openAndReadFile(path string) []byte {
	// open file
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("unable to open test file %s: %s", path, err))
	}
	// read file
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("unable to ReadAll of test file %s: %s", path, err))
	}
	return b
}
