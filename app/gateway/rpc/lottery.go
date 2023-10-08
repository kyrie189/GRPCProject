package rpc

import (
	"context"
	pb "github.com/CocaineCong/grpc-todolist/idl/pb/lottery"
)

func InitAward(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	_, err := LotteryClient.InitAward(ctx,in)
	if err != nil{
		return nil, err
	}
	return nil,err
}
func Draw(ctx context.Context, in *pb.DrawRequest) (*pb.DrawResponse, error) {
	out, err := LotteryClient.Draw(ctx,in)
	return out,err
}
func ListAwardInfo(ctx context.Context, in *pb.DrawRequest) (*pb.ListAwardInfoResponse, error) {
	out, err := LotteryClient.ListAwardInfo(ctx,in)
	return out,err
}
func ToMysql(ctx context.Context, in *pb.Empty) (*pb.ToMysqlResponse, error){
	out, err := LotteryClient.ToMysql(ctx,in)
	return out,err
}