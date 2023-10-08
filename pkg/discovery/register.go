package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *logrus.Logger
}
// NewRegister create a register based on etcd
func NewRegister(etcdAddrs []string, logger *logrus.Logger) *Register {



	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register a user
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}

	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{}) //

	go r.keepAlive()

	return r.closeCh, nil
}
// 正式的etcd注册
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()
	// 超时上下文防止授予租约超时
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)   // 授予租约
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID
	// keepAlive每次续租，都会在通道中塞入数据
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo) // 将结构体Marshal成一个json之后作为一个value
	if err != nil {
		return err
	}
	// BuildRegisterPath ==> 服务名+addr
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))

	return err
}

// Stop stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)

	for {
		select {
		case <-r.closeCh:   // 关闭
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed, error: ", err)
			}

			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed, error: ", err)
			}
		case res := <-r.keepAliveCh:
			if res == nil {   // 续租失败
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error: ", err)
				}
			}
			// fmt.Printf("service %s(%s) keep alive\n",r.srvInfo.Name,r.srvInfo.Addr)
		case <-ticker.C: // 超时
			if r.keepAliveCh == nil {  // 因为续租超时
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error: ", err)
				}
			}
		}
	}
}

func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		weightstr := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(weightstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}

			_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write([]byte("update server weight success"))
	})
}

func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegisterPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}

	server := Server{}
	fmt.Println(len(resp.Kvs))
	// 添加实现负载均衡
	//if len(resp.Kvs)==0{
	//	return server,errors.New("无服务")
	//}
	//// 随机选择一个服务
	//rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	//randIndex :=rd.Intn(len(resp.Kvs)) // [0,n)
	//
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &server); err != nil {
			return server, err
		}
	}

	return server, err
}
