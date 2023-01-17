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

type FileRenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileRenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileRenameLogic {
	return &FileRenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileRenameLogic) FileRename(req *types.FileRenameRequest, uid string) (resp *types.FileRenameReply, err error) {
	resp = &types.FileRenameReply{Code: 0}

	// 参数检验
	if req.Fid == "" {
		resp.Code = config.FILE_RID_EMPTY
		return
	}
	if req.Name == "" {
		resp.Code = config.FOLDER_NAME_EMPTY
		return
	}

	// 检查用户云盘中该文件名是否已存在
	total, err := l.svcCtx.Engine.Where("name = ? and uid = ?", req.Name, uid).Count(new(models.UserFileBasic))
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if total > 0 {
		resp.Code = config.FILE_NAME_HAS
		return
	}

	// 检查文件是否存在和文件操作权限
	userFile := new(models.UserFileBasic)
	has, err := l.svcCtx.Engine.Where("fid = ?", req.Fid).Select("uid").Get(userFile)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = config.FILE_NOT_EXIST
		return
	}
	if userFile.Uid != uid {
		resp.Code = config.ACCESS_DENIED
		return
	}

	// 修改文件名
	num, err := l.svcCtx.Engine.Where("fid = ?", req.Fid).Cols("name").Update(&models.UserFileBasic{Name: req.Name})
	if err != nil || num < 1 {
		if num < 1 {
			resp.Code = config.SERVER_PANIC
			return
		}
		return nil, err
	}
	return
}
