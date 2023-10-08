package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"

	"github.com/CocaineCong/grpc-todolist/app/gateway/rpc"
	pb "github.com/CocaineCong/grpc-todolist/idl/pb/task"
	"github.com/CocaineCong/grpc-todolist/pkg/ctl"
)
// api
func GetTaskList(ctx *gin.Context) {
	var req pb.TaskRequest

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func CreateTask(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskCreate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func UpdateTask(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskUpdate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func DeleteTask(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserID = user.Id
	r, err := rpc.TaskDelete(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// 上传文件
func UploadFile(c *gin.Context){
	// 从请求中获取上传的文件
	form, err := c.MultipartForm()  // multipartForm
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 获取所有上传的文件
	files := form.File["files"]


	// 逐块发送文件内容
	buffer := make([]byte, 1024) // 根据实际需求选择适当的缓冲区大小
	for _,file := range files{
		// 创建流式上传文件的客户端
		stream, err := rpc.TaskClient.UploadFile(context.Background())
		src,err := file.Open()
		defer src.Close()
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  // 文件打开失败
		}

		for {
			n, err := src.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				log.Fatalf("Failed to read file: %v", err)
			}

			// 发送文件块到服务端
			err = stream.Send(&pb.FileInfo{
				Content: buffer[:n],
				Name: file.Filename,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				log.Fatalf("Failed to send file chunk: %v", err)
			}
		}
		_, err = stream.CloseAndRecv()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Fatalf("Failed to close stream : %v", err)
		}
	}


	c.JSON(http.StatusOK, gin.H{"上传成功":"成功"})  // 文件打开失败
}