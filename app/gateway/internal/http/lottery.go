package http

import (
	"github.com/CocaineCong/grpc-todolist/app/gateway/rpc"
	pb "github.com/CocaineCong/grpc-todolist/idl/pb/lottery"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/grpc-todolist/pkg/ctl"
)

// 初始化数据
func InitAward(ctx *gin.Context) {
	req := pb.Empty{}
	_,err := rpc.InitAward(ctx,&req)
	if err != nil{
		ctx.JSON(http.StatusCreated, ctl.RespSuccess(ctx, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, "创建奖品成功"))
}
// 抽奖
func Draw(ctx *gin.Context) {
	req := pb.DrawRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	//fmt.Println("我已经拿到数据啦",req.Id)
	out,err := rpc.Draw(ctx,&req)
	if err != nil{
		ctx.JSON(http.StatusCreated, ctl.RespError(ctx, err,"抽奖失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, out))
}
func ListAwardInfo(ctx *gin.Context) {
	req := pb.DrawRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	out,err := rpc.ListAwardInfo(ctx,&req)
	if err != nil{
		ctx.JSON(http.StatusCreated, ctl.RespError(ctx, err,"抽奖失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, out))

}
func ToMysql(ctx *gin.Context) {
	req := pb.Empty{}

	out,err := rpc.ToMysql(ctx,&req)
	if err != nil{
		ctx.JSON(http.StatusCreated, ctl.RespError(ctx, err,"加入mysql失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, out))

}