package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/grpc-todolist/app/user/internal/repository/db/dao"
	"github.com/CocaineCong/grpc-todolist/app/user/internal/service"
	"github.com/CocaineCong/grpc-todolist/config"
	pb "github.com/CocaineCong/grpc-todolist/idl/pb/user"
	"github.com/CocaineCong/grpc-todolist/pkg/discovery"
)

func main() {
	config.InitConfig()
	dao.InitDB()

	//etcd
	etcdAddress := []string{config.Conf.Etcd.Address}	// etcd 地址、
	// 多个服务地址，循环添加到etcd的节点中
	for _,v:= range config.Conf.Services["user"].Addr{
		etcdRegister := discovery.NewRegister(etcdAddress, logrus.New()) // etcd服务注册
		defer etcdRegister.Stop()
		grpcAddress := v   // 服务地址
		fmt.Println(v)
		// etcd的服务节点
		userNode := discovery.Server{
			Name: config.Conf.Domain["user"].Name,
			Addr: grpcAddress,
		}
		// grpc服务
		server := grpc.NewServer()
		defer server.Stop()
		// 绑定service
		pb.RegisterUserServiceServer(server, service.GetUserSrv())
		lis, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			panic(any(err))
		}
		// 绑定服务之后注册到etcd
		if _, err := etcdRegister.Register(userNode, 10); err != nil {
			panic(any(fmt.Sprintf("start server failed, err: %v", err)))
		}
		logrus.Info("server started listen on ", grpcAddress)
		// 开启协程进行监听grpc服务
		go func(){
			if err := server.Serve(lis); err != nil {
				panic(any(err))
			}
		}()

	}
	fmt.Println("服务启动")

	select {

	}
}
