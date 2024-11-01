package main

import (
	"flag"
	"net"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/CavnHan/wallet-chain-account/chaindispatcher"
	"github.com/CavnHan/wallet-chain-account/config"
	"github.com/CavnHan/wallet-chain-account/rpc/account"
)

func main() {
	//定义config并解析
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()

	//创建config
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}

	//创建调度器
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	//创建grpc服务
	server := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer server.GracefulStop()

	//注册服务
	account.RegisterWalletAccountServiceServer(server, dispatcher)

	//监听端口
	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}

	//注册反射服务
	reflection.Register(server)

	log.Info("wallet rpc services start success", "port:", conf.Server.Port)

	//启动服务
	if err := server.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err:", err)
		panic(err)
	}
}
