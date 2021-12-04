module clientv2_6_0

go 1.13

require (
	github.com/vesoft-inc/nebula-go/v2 v2.6.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/logger v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/types v0.0.0
)

replace (
	github.com/vesoft-inc/nebula-http-gateway/pkg/logger v0.0.0 => ../../logger
	github.com/vesoft-inc/nebula-http-gateway/pkg/types v0.0.0 => ../../types
)
