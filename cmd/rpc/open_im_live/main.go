package main

import (
	"Open_IM/internal/rpc/live"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImLivePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.UserPrometheusPort[0], "userPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start user rpc live, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion)
	rpcServer := live.NewRpcLiveServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
