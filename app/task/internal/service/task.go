package service

import (
	"context"
	"github.com/CocaineCong/grpc-todolist/pkg/storage"
	"io"
	"sync"

	"github.com/CocaineCong/grpc-todolist/app/task/internal/repository/db/dao"
	pb "github.com/CocaineCong/grpc-todolist/idl/pb/task"
	"github.com/CocaineCong/grpc-todolist/pkg/e"
)

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
	pb.UnimplementedTaskServiceServer
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})
	return TaskSrvIns
}
func (*TaskSrv) TaskCreate(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskCommonResponse, err error) {
	resp = new(pb.TaskCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).CreateTask(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (*TaskSrv) TaskShow(ctx context.Context, req *pb.TaskRequest) (resp *pb.TasksDetailResponse, err error) {
	resp = new(pb.TasksDetailResponse)
	r, err := dao.NewTaskDao(ctx).ListTaskByUserId(req.UserID)
	resp.Code = e.SUCCESS
	if err != nil {
		resp.Code = e.ERROR
		return
	}
	for i := range r {
		resp.TaskDetail = append(resp.TaskDetail, &pb.TaskModel{
			TaskID:    r[i].TaskID,
			UserID:    r[i].UserID,
			Status:    int64(r[i].Status),
			Title:     r[i].Title,
			Content:   r[i].Content,
			StartTime: r[i].StartTime,
			EndTime:   r[i].EndTime,
		})
	}
	return
}

func (*TaskSrv) TaskUpdate(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskCommonResponse, err error) {
	resp = new(pb.TaskCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (*TaskSrv) TaskDelete(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskCommonResponse, err error) {
	resp = new(pb.TaskCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).DeleteTaskById(req.TaskID, req.UserID)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

// 上传文件
func (*TaskSrv) UploadFile(stream pb.TaskService_UploadFileServer) error{
	file := storage.NewFile("")
	for {
		// 从文件流中接收文件
		req, err := stream.Recv()
		if err == io.EOF {  // 读取结束
			// 客户端调用CloseAndRecv时
			if err := storage.Store(file);err != nil{
				return err
			}
			// 文件传输完成
			return stream.SendAndClose(&pb.UploadFileResponse{})  // 返回
		}
		if err != nil {
			return err
		}
		if file.Name == ""{
			file.Name = req.Name  // 设定文件名
		}
		// 将数据流写入buffer
		if err := file.Write(req.GetContent()); err != nil {
			return err
		}
	}
}