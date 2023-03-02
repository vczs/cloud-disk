package logic

import (
	"context"
	"database/sql"

	"cloud-disk/disk/define"
	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileShareLogic {
	return &FileShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileShareLogic) FileShare(req *types.FileShareRequest, uid string) (resp *types.FileShareReply, err error) {
	resp = &types.FileShareReply{Code: 0}

	// 参数校验
	if req.Fid == "" {
		resp.Code = define.FILE_RID_EMPTY
		return
	}

	// 检查要操作的文件是否存在
	userFile := new(models.UserFileBasic)
	has, err := l.svcCtx.Engine.Where("fid = ? and uid = ?", req.Fid, uid).Cols("type", "rid").Get(userFile)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = define.FILE_NOT_EXIST
		return
	}

	// 检查要操作的文件是否为文件
	if userFile.Type != define.FileTypeFile {
		resp.Code = define.FOLDER_NOT_SHARE
		return
	}

	// 构建数据
	sid := helper.GetUid()
	share := &models.ShareBasic{
		Sid:     sid,
		Uid:     uid,
		Fid:     req.Fid,
		Type:    define.ShareTypePublic,
		Encrypt: sql.NullString{Valid: false},
		Expired: req.Expired,
		Browse:  0,
	}
	if req.Encrypt != "" {
		share.Type = define.ShareTypeEncrypt
		share.Encrypt = sql.NullString{String: req.Encrypt, Valid: true}
	}

	// 分享文件
	affect, err := l.svcCtx.Engine.Insert(share)
	if err != nil {
		return nil, err
	}
	if affect < 1 {
		resp.Code = define.SHARE_FAILED
	}

	resp.Data.Sid = sid
	return
}
