package rpc

import (
	"context"
	"fmt"
	"github.com/CocaineCong/grpc-todolist/idl/pb/lottery"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/CocaineCong/grpc-todolist/config"
	"github.com/CocaineCong/grpc-todolist/idl/pb/task"
	"github.com/CocaineCong/grpc-todolist/idl/pb/user"
	"github.com/CocaineCong/grpc-todolist/pkg/discovery"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	UserClient user.UserServiceClient
	TaskClient task.TaskServiceClient
	LotteryClient lottery.LotteryServiceClient
)

func  Init() {
	// 连接etcd
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	initClient(config.Conf.Domain["user"].Name, &UserClient)
	initClient(config.Conf.Domain["task"].Name, &TaskClient)
	initClient(config.Conf.Domain["lottery"].Name, &LotteryClient)
}

func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(any(err))
	}

	switch c := client.(type) {
	case *user.UserServiceClient:
		*c = user.NewUserServiceClient(conn)
	case *task.TaskServiceClient:
		*c = task.NewTaskServiceClient(conn)
	case *lottery.LotteryServiceClient:
		*c = lottery.NewLotteryServiceClient(conn)
	default:
		panic(any("unsupported client type"))
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	// etcd://user
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)
	//fmt.Printf("service :%s\taddr: %s\n",serviceName,addr)
	// Load balance
	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}
	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
