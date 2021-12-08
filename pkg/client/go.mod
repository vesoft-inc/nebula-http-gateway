module github.com/vesoft-inc/nebula-http-gateway/pkg/client

go 1.13

require (
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_0_0_ga v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_5_0 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_5_1 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_6_0 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/logger v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/types v0.0.0
)

replace (
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_0_0_ga v0.0.0 => ./clientv2_0_0_ga
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_5_0 v0.0.0 => ./clientv2_5_0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_5_1 v0.0.0 => ./clientv2_5_1
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/clientv2_6_0 v0.0.0 => ./clientv2_6_0
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/logger v0.0.0 => ./logger
	github.com/vesoft-inc/nebula-http-gateway/pkg/client/types v0.0.0 => ./types
)
