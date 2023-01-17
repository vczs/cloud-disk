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

type FolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FolderCreateLogic {
	return &FolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FolderCreateLogic) FolderCreate(req *types.FolderCreateRequest, uid string) (resp *types.FolderCreateReply, err error) {
	resp = &types.FolderCreateReply{Code: 0}

	// 参数检验
	if req.Name == "" {
		resp.Code = config.FOLDER_NAME_EMPTY
		return
	}

	number, err := l.svcCtx.Engine.Where("uid = ? and type = ?", uid, config.FileTypeFolder).Count(new(models.UserFileBasic))
	if err != nil {
		return
	}
	if number > int64(config.MaxUserFolder) {
		resp.Code = config.FOLDER_NUMBER_MAX
		return
	}

	// 默认查询一级文件夹列表
	pid := req.Pid
	if pid == "" || pid == "0" {
		pid = "0"
	} else {
		// 检查文件是否存在和文件操作权限
		uf := new(models.UserFileBasic)
		has, err := l.svcCtx.Engine.Where("fid = ? and type = ?", pid, config.FileTypeFolder).Select("uid").Get(uf)
		if err != nil {
			return nil, err
		}
		if !has {
			resp.Code = config.FOLDER_NOT_EXIST
			return resp, nil
		}
		if uf.Uid != uid {
			resp.Code = config.ACCESS_DENIED
			return resp, nil
		}
	}

	// 检查重名
	num, err := l.svcCtx.Engine.Where("pid = ? and name = ? and uid = ?", pid, req.Name, uid).Count(new(models.UserFileBasic))
	if err != nil {
		return nil, err
	}
	if num > 0 {
		resp.Code = config.FILE_NAME_HAS
		return
	}

	// 构造数据
	fid := helper.GetUid()
	userFile := &models.UserFileBasic{
		Fid:    fid,
		Uid:    uid,
		Pid:    pid,
		Type:   config.FileTypeFolder,
		Name:   req.Name,
		Number: 0,
	}

	userFileAffect, err := l.svcCtx.Engine.Insert(userFile)
	if err != nil || userFileAffect < 1 {
		panic(config.SERVER_PANIC)
	}

	resp.Data.Fid = fid
	return
}
