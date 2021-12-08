module github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_0_0_ga

go 1.13

require (
	github.com/vesoft-inc/nebula-go/v2 v2.0.0-ga
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/logger v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/types v0.0.0
)

replace (
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/logger v0.0.0 => ../logger
    github.com/vesoft-inc/nebula-http-gateway/pkg/client/types v0.0.0 => ../types
)
