module client

go 1.13

require (
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_0_0_ga v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_0 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_1 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_6_0 v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/logger v0.0.0
	github.com/vesoft-inc/nebula-http-gateway/pkg/types v0.0.0
)

replace (
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_0_0_ga v0.0.0 => ./../adapters/clientv2_0_0_ga
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_0 v0.0.0 => ./../adapters/clientv2_5_0
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_1 v0.0.0 => ./../adapters/clientv2_5_1
	github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_6_0 v0.0.0 => ../adapters/clientv2_6_0
	github.com/vesoft-inc/nebula-http-gateway/pkg/logger v0.0.0 => ../logger
	github.com/vesoft-inc/nebula-http-gateway/pkg/types v0.0.0 => ../types
)
