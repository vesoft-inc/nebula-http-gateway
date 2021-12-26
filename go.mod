module github.com/vesoft-inc/nebula-http-gateway

go 1.13

replace github.com/vesoft-inc/nebula-http-gateway/ccore v0.0.0 => ./ccore

require (
	github.com/astaxie/beego v1.12.3
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/vesoft-inc/nebula-importer v1.0.1-0.20211213064541-05a8646be295
)

require (
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/vesoft-inc/nebula-go/v2 v2.5.2-0.20211221081231-40030d441885 // indirect
	github.com/vesoft-inc/nebula-http-gateway/ccore v0.0.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210105210732-16f7687f5001 // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
