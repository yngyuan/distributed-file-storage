package config

import (
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

const (
	// UploadServiceHost : 上传服务监听的地址
	UploadServiceHost = "0.0.0.0:8080"
	// UploadLBHost : 上传服务LB地址
	UploadLBHost = "http://upload.fileserver.com"
	// DownloadLBHost : 下载服务LB地址
	DownloadLBHost = "http://download.fileserver.com"
	// TracerAgentHost : tracing agent地址
	TracerAgentHost = "127.0.0.1:6831"
)

// RegistryConsul : 配置 consul
func RegistryConsul() registry.Registry {
	return consul.NewRegistry(
		// TODO ip需根据实际情况来修改
		registry.Addrs("192.168.2.244:8500"),
	)
}

// RegistryClient : 注册中心client
func RegistryClient(r registry.Registry) selector.Selector {
	return selector.NewSelector(
		selector.Registry(r),                      //传入consul注册
		selector.SetStrategy(selector.RoundRobin), //指定查询机制
	)
}
