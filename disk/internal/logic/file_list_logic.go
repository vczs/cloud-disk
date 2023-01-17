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

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListRequest, uid string) (resp *types.FileListReply, err error) {
	resp = &types.FileListReply{Code: 0}

	// 默认查询一级文件夹列表
	fid := req.Fid
	if fid == "" {
		fid = "0"
	}
	// 检查目标文件夹是否存在
	pid := "0"
	if fid != "0" {
		uf := new(models.UserFileBasic)
		has, err := l.svcCtx.Engine.Where("fid = ? and uid = ?", fid, uid).Cols("type", "pid").Get(uf)
		if err != nil {
			helper.VczsLog("", err)
			return nil, err
		}
		if !has {
			resp.Code = config.FOLDER_NOT_EXIST
			return resp, nil
		}
		// 检查目标文件是否为文件夹
		if uf.Type != config.FileTypeFolder {
			resp.Code = config.TARGET_NOT_FOLDER
			return resp, nil
		}
		pid = uf.Pid
	}

	// 处理分页信息
	number := req.Size
	if number == 0 {
		number = config.DefaultPageSize
	}
	start := req.Page
	if start == 0 {
		start = 1
	}
	start = (start - 1) * number

	// 分页查找
	files := make([]*types.Files, 0)
	total, err := l.svcCtx.Engine.Table(new(models.UserFileBasic).TableName()).Alias("uf").
		Where("uf.uid = ? and uf.pid = ?", uid, fid).Select("uf.fid, uf.name, b.ext, uf.type, b.url, uf.number, uf.size").
		Join("LEFT", []string{new(models.BucketBasic).TableName(), "b"}, "uf.rid = b.rid").
		// 这里通过[]string给 new(models.BucketBasic).TableName()获取的表 起了一个别名"b"
		Where("uf.dt IS NULL").Limit(number, start).FindAndCount(&files)
	if err != nil {
		return nil, err
	}

	// 返回数据
	resp.Data.Pid = pid
	resp.Data.List = files
	resp.Data.Total = total
	return
}
