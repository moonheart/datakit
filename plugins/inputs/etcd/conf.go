package etcd

import "time"

const (
	pluginName = "etcd"

	etcdConfigSample = `
# [etcd]
# [[etcd.subscribes]]
#       ## etcd 地址
#	host = "127.0.0.1"
#
#       ## etcd 端口
#	port = "2379"
#
#	## 采集周期，时间单位是秒
#	collect_cycle = 60
#
#       ## measurement，不可重复
#       measurement = "etcd"
`
)

type Subscribe struct {
	EtcdHost    string        `toml:"host"`
	EtcdPort    int           `toml:"port"`
	Cycle       time.Duration `toml:"collect_cycle"`
	Measurement string        `toml:"measurement"`
}

type Config struct {
	Subscribes []Subscribe `toml:"subscribes"`
}
