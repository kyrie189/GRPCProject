package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CocaineCong/grpc-todolist/app/lottery/internal/repository/db/dao"
	"github.com/CocaineCong/grpc-todolist/app/lottery/internal/repository/db/model"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"sync"
	"time"

	pb "github.com/CocaineCong/grpc-todolist/idl/pb/lottery"
)

var LotterySrvIns *LotterySrv
var LotterySrvOnce sync.Once

type LotterySrv struct {
	pb.UnimplementedLotteryServiceServer
}
// 工厂模式获取实例
func GetLotterySrv() *LotterySrv{
	LotterySrvOnce.Do(func(){
		LotterySrvIns = &LotterySrv{}
	})
	return LotterySrvIns
}

// 初始化redis数据库
func (l *LotterySrv) InitAward(c context.Context, req *pb.Empty) (resp *pb.Empty, err error) {
	resp = new(pb.Empty)
	rb := dao.NewRedisClient()  // 获取redis客户端

	err = rb.Set(context.Background(),"phone",200,0).Err()
	err = rb.Set(context.Background(),"phoneShell",2000,0).Err()
	err = rb.Set(context.Background(),"20RMB",5000,0).Err()
	err = rb.Set(context.Background(),"5RMB",20000,0).Err()

	return
}
func (l *LotterySrv) Draw(ctx context.Context, in *pb.DrawRequest) (out *pb.DrawResponse, err error) {
	out = new(pb.DrawResponse)
	out.AwardInfo = new(pb.AwardInfo)
	//获取随机数并抽奖
	giftNumber := 5
	code := luckyCode(int64(in.Id)) 	// 获取随机数
	//fmt.Println("获取随机数: ", code)
	if code >= 1 && code <= 20 {
		giftNumber = 1
	} else if code > 20 && code <= 220 {  //
		giftNumber = 2
	} else if code > 220 && code <= 2200 {
		giftNumber = 3
	} else if code > 2200 && code <= 7200 {
		giftNumber = 4
	}
	if giftNumber == 5{
		out.Msg = "no win"
		return out,nil
	}

	switch giftNumber {
	case 1:
		out,err = GetAwardInfo(ctx,"phone",in.Id,out)
	case 2:
		out,err = GetAwardInfo(ctx,"phoneShell",in.Id,out)
	case 3:
		out,err = GetAwardInfo(ctx,"20RMB",in.Id,out)
	case 4:
		out,err = GetAwardInfo(ctx,"5RMB",in.Id,out)
	}

	return out,nil
}
func GetAwardInfo(ctx context.Context,award string,id int32,out *pb.DrawResponse)(*pb.DrawResponse, error) {

	rdb := dao.NewRedisClient()
	var err error

	//乐观锁保证事务的一致性
	// 定义一个回调函数，用于处理事务逻辑
	fn := func(tx *redis.Tx) error {
		// 先查询下当前watch监听的key的值
		val, err := tx.Get(ctx, award).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		// 这里可以处理业务
		if val <= 0 {
			out.Msg = "sorry,没有中奖"
			return err
		}
		// 用户中奖信息
		out.AwardInfo.Created = time.Now().Unix()
		out.AwardInfo.OrderId = id
		out.AwardInfo.Award = award
		out.Msg =  fmt.Sprintf("恭喜中奖：%s",award)
		result ,err := json.Marshal(out.AwardInfo)
		if err != nil{
			return err
		}
		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在这里给key设置最新值
			err := pipe.Decr(ctx,award).Err()  // 加入队列
			if err != nil{
				return err
			}
			// 添加用户中奖信  （怎么保证加入失败的问题）
			if err := pipe.RPush(ctx,"user_"+strconv.Itoa(int(id)),result).Err();err != nil{
				return err
			}
			return nil
		})
		return err
	}

	// 使用Watch监听一些Key, 同时绑定一个回调函数fn, 监听Key后的逻辑写在fn这个回调函数里面
	// 如果想监听多个key，可以这么写：client.Watch(ctx,fn, "key1", "key2", "key3")
	err = rdb.Watch(ctx, fn, award)

	return out,err
}

func luckyCode(id int64) int32 {
	rateMax := 10000
	seed := time.Now().UnixNano()
	seed = seed + id
	r := rand.New(rand.NewSource(seed))
	code := r.Int31n(int32(rateMax)) + 1
	return code //返回一个随机数[1,10000]
}


// 列出中奖信息
func (l *LotterySrv) ListAwardInfo(ctx context.Context, in *pb.DrawRequest) (out *pb.ListAwardInfoResponse, err error){
	out = new(pb.ListAwardInfoResponse)
	rdb := dao.NewRedisClient()
	// 返回从0开始到-1位置之间的数据，意思就是返回全部数据
	vals, err := rdb.LRange(ctx,"user_"+strconv.Itoa(int(in.Id)),0,-1).Result()
	if err != nil {
		return nil,err
	}
	for _,v := range vals{
		temp := pb.AwardInfo{}
		if err = json.Unmarshal([]byte(v),&temp);err != nil{
			return nil,err
		}
		out.List = append(out.List,&temp)
	}
	return out ,nil
}

//存入mysql
func  (l *LotterySrv)ToMysql(ctx context.Context, in *pb.Empty) (out *pb.ToMysqlResponse, err error){
	out = new(pb.ToMysqlResponse)
	rdb := dao.NewRedisClient()
	db := dao.NewDBClient(ctx)
	name := "user_"
	//rdbTemp := pb.AwardInfo{}


	for i:=0;i<100000;i++{

		key := name + strconv.Itoa(i)
		//fmt.Println(key)
		val,err := rdb.LRange(context.Background(),key,0,-1).Result()
		if err != nil && err != redis.Nil {
			return out,nil
		}

		// 遍历user_1的中奖信息
		for _,v := range val{
			dbTemp := model.Lottery{}
			if err = json.Unmarshal([]byte(v),&dbTemp);err != nil{
				fmt.Println("解析失败：" ,v)
				continue
				// return nil,err
			}
			//dbTemp.UserId = rdbTemp.OrderId
			//dbTemp.Award = rdbTemp.Award
			//dbTemp.LuckyTime = rdbTemp.Created
			if err = db.Create(&dbTemp).Error;err != nil {
				fmt.Println("加入mysql失败：" ,v)
			}
		}
	}

	return out ,err
}