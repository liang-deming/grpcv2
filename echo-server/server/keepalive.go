package server

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func GetKeepaliveOpt() (opts []grpc.ServerOption) {
	//服务端轻质保活策略，客户端违反该策略将被关闭
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, //客户端空闲超时时间
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		//客户端空闲5秒，发送ping保活
		Time: 5 * time.Second,
		// ping 超时时间
		Timeout: 1 * time.Second,
	}

	return []grpc.ServerOption{grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp)}
}
