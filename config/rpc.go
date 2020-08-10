package config

import (
        "time"
        "github.com/micro/go-micro/client"
)

// RpcOpts 设置RPC请求参数　(超时时间等)
var RpcOpts client.CallOption = func(o *client.CallOptions) {
    o.RequestTimeout = time.Second * 30
    o.DialTimeout = time.Second * 30
}
