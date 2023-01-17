package logic

import (
	"context"

	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileMoveLogic {
	return &FileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileMoveLogic) FileMove(req *types.FileMoveRequest, uid string) (resp *types.FileMoveReply, err error) {
	resp = &types.FileMoveReply{Code: 0}

	// 参数校验
	if req.Fid == "" {
		resp.Code = config.FILE_RID_EMPTY
		return
	}
	if req.Pid == "" {
		resp.Code = config.FILE_PID_EMPTY
		return
	}

	// 检查要操作的文件是否存在
	userFile := new(models.UserFileBasic)
	has, err := l.svcCtx.Engine.Where("fid = ? and uid = ?", req.Fid, uid).Select("name,size,pid").Get(userFile)
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if !has {
		resp.Code = config.FILE_NOT_EXIST
		return
	}
	if req.Pid == userFile.Pid {
		return
	}

	// 检查目标文件是否存在
	pid := req.Pid
	if pid != "0" {
		uf := new(models.UserFileBasic)
		has, err = l.svcCtx.Engine.Where("fid = ? and uid = ?", pid, uid).Select("type,size").Get(uf)
		if err != nil {
			helper.VczsLog("", err)
			return nil, err
		}
		if !has {
			resp.Code = config.FOLDER_NOT_EXIST
			return
		}
		// 检查目标文件是否为文件夹
		if uf.Type != config.FileTypeFolder {
			resp.Code = config.TARGET_NOT_FOLDER
			return
		}
	}

	// 检查目标文件夹下是否有与要操作的文件 重名的文件
	num, err := l.svcCtx.Engine.Where("pid = ? and name = ? and uid = ?", pid, userFile.Name, uid).Count(new(models.UserFileBasic))
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if num > 0 {
		resp.Code = config.TARGET_PATH_NAME_HAS
		return
	}

	// 创建事务
	session := l.svcCtx.Engine.NewSession()
	defer session.Close()
	// 开启事务
	if err = session.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		err := recover()
		if err != nil {
			session.Rollback()
			resp.Code = config.FILE_MOVE_FAILED
		} else {
			session.Commit()
		}
	}()

	// 移动文件
	if pid != "0" {
		err = FolderSizeChange(session, config.FolderSizeChangeIncrease, pid, userFile.Size, 1)
		if err != nil {
			panic(config.FILE_MOVE_FAILED)
		}
	}
	if userFile.Pid != "0" {
		err = FolderSizeChange(session, config.FolderSizeChangeDecrease, userFile.Pid, userFile.Size, 1)
		if err != nil {
			panic(config.FILE_MOVE_FAILED)
		}
	}
	affect, err := session.Where("fid = ?", req.Fid).Cols("pid").Update(&models.UserFileBasic{Pid: pid})
	if err != nil || affect < 1 {
		panic(config.FILE_MOVE_FAILED)
	}

	return
}
