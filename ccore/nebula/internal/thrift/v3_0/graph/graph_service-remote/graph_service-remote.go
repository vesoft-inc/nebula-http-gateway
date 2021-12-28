// Autogenerated by Thrift Compiler (facebook)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
// @generated

package main

import (
	"../../github.com/vesoft-inc/nebula-go/v2/nebula/graph"
	"flag"
	"fmt"
	thrift "github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  AuthResponse authenticate(string username, string password)")
	fmt.Fprintln(os.Stderr, "  void signout(i64 sessionId)")
	fmt.Fprintln(os.Stderr, "  ExecutionResponse execute(i64 sessionId, string stmt)")
	fmt.Fprintln(os.Stderr, "  ExecutionResponse executeWithParameter(i64 sessionId, string stmt,  parameterMap)")
	fmt.Fprintln(os.Stderr, "  string executeJson(i64 sessionId, string stmt)")
	fmt.Fprintln(os.Stderr, "  string executeJsonWithParameter(i64 sessionId, string stmt,  parameterMap)")
	fmt.Fprintln(os.Stderr, "  VerifyClientVersionResp verifyClientVersion(VerifyClientVersionReq req)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	var parsedUrl url.URL
	var trans thrift.Transport
	_ = strconv.Atoi
	_ = math.Abs
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Parse()

	if len(urlString) > 0 {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewHTTPPostClient(parsedUrl.String())
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewSocket(thrift.SocketAddr(net.JoinHostPort(host, portStr)))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.ProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	client := graph.NewGraphServiceClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "authenticate":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "Authenticate requires 2 args")
			flag.Usage()
		}
		argvalue0 := []byte(flag.Arg(1))
		value0 := argvalue0
		argvalue1 := []byte(flag.Arg(2))
		value1 := argvalue1
		fmt.Print(client.Authenticate(value0, value1))
		fmt.Print("\n")
		break
	case "signout":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Signout requires 1 args")
			flag.Usage()
		}
		argvalue0, err17 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err17 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Signout(value0))
		fmt.Print("\n")
		break
	case "execute":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "Execute requires 2 args")
			flag.Usage()
		}
		argvalue0, err18 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err18 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		argvalue1 := []byte(flag.Arg(2))
		value1 := argvalue1
		fmt.Print(client.Execute(value0, value1))
		fmt.Print("\n")
		break
	case "executeWithParameter":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "ExecuteWithParameter requires 3 args")
			flag.Usage()
		}
		argvalue0, err20 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err20 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		argvalue1 := []byte(flag.Arg(2))
		value1 := argvalue1
		arg22 := flag.Arg(3)
		mbTrans23 := thrift.NewMemoryBufferLen(len(arg22))
		defer mbTrans23.Close()
		_, err24 := mbTrans23.WriteString(arg22)
		if err24 != nil {
			Usage()
			return
		}
		factory25 := thrift.NewSimpleJSONProtocolFactory()
		jsProt26 := factory25.GetProtocol(mbTrans23)
		containerStruct2 := graph.NewGraphServiceExecuteWithParameterArgs()
		err27 := containerStruct2.ReadField3(jsProt26)
		if err27 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.ParameterMap
		value2 := argvalue2
		fmt.Print(client.ExecuteWithParameter(value0, value1, value2))
		fmt.Print("\n")
		break
	case "executeJson":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "ExecuteJson requires 2 args")
			flag.Usage()
		}
		argvalue0, err28 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err28 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		argvalue1 := []byte(flag.Arg(2))
		value1 := argvalue1
		fmt.Print(client.ExecuteJson(value0, value1))
		fmt.Print("\n")
		break
	case "executeJsonWithParameter":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "ExecuteJsonWithParameter requires 3 args")
			flag.Usage()
		}
		argvalue0, err30 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err30 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		argvalue1 := []byte(flag.Arg(2))
		value1 := argvalue1
		arg32 := flag.Arg(3)
		mbTrans33 := thrift.NewMemoryBufferLen(len(arg32))
		defer mbTrans33.Close()
		_, err34 := mbTrans33.WriteString(arg32)
		if err34 != nil {
			Usage()
			return
		}
		factory35 := thrift.NewSimpleJSONProtocolFactory()
		jsProt36 := factory35.GetProtocol(mbTrans33)
		containerStruct2 := graph.NewGraphServiceExecuteJsonWithParameterArgs()
		err37 := containerStruct2.ReadField3(jsProt36)
		if err37 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.ParameterMap
		value2 := argvalue2
		fmt.Print(client.ExecuteJsonWithParameter(value0, value1, value2))
		fmt.Print("\n")
		break
	case "verifyClientVersion":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "VerifyClientVersion requires 1 args")
			flag.Usage()
		}
		arg38 := flag.Arg(1)
		mbTrans39 := thrift.NewMemoryBufferLen(len(arg38))
		defer mbTrans39.Close()
		_, err40 := mbTrans39.WriteString(arg38)
		if err40 != nil {
			Usage()
			return
		}
		factory41 := thrift.NewSimpleJSONProtocolFactory()
		jsProt42 := factory41.GetProtocol(mbTrans39)
		argvalue0 := graph.NewVerifyClientVersionReq()
		err43 := argvalue0.Read(jsProt42)
		if err43 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.VerifyClientVersion(value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
