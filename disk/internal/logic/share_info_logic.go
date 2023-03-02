package logic

import (
	"context"
	"time"

	"cloud-disk/disk/define"
	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareInfoLogic {
	return &ShareInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareInfoLogic) ShareInfo(req *types.ShareInfoRequest) (resp *types.ShareInfoReply, err error) {
	resp = &types.ShareInfoReply{Code: 0}

	// 参数校验
	if req.Sid == "" {
		resp.Code = define.FILE_SID_EMPTY
		return
	}

	// 检查分享资源
	share := new(models.ShareBasic)
	has, err := l.svcCtx.Engine.Where("sid = ?", req.Sid).Cols("fid", "type", "encrypt", "expired").Get(share)
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if !has {
		resp.Code = define.SHARE_NOT_EXIST
		return
	} else {
		userFile := new(models.UserFileBasic)
		num, err := l.svcCtx.Engine.Where("fid = ?", share.Fid).Count(userFile)
		if err != nil {
			helper.VczsLog("", err)
			return nil, err
		}
		if num < 1 {
			resp.Code = define.SHARE_NOT_EXIST
			return resp, nil
		}
	}

	// 检查资源是否在有效期
	if share.Expired < time.Now().Unix() {
		resp.Code = define.SHARE_EXPIRE
		return
	}

	// 检查加密资源密码
	if share.Type == define.ShareTypeEncrypt {
		if req.Encrypt == "" {
			resp.Code = define.SHARE_EMPTY
			return
		} else {
			if req.Encrypt != share.Encrypt.String {
				resp.Code = define.SHARE_PWD_ERR
				return
			}
		}
	}

	// 对分享记录的点击次数进行 + 1
	_, err = l.svcCtx.Engine.Exec("UPDATE "+new(models.ShareBasic).TableName()+" SET browse = browse + 1 WHERE sid = ?", req.Sid)
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}

	// 查询资源
	data := types.ShareInfoReplyData{}
	has, err = l.svcCtx.Engine.Table(new(models.ShareBasic).TableName()).Alias("s").
		Where("s.sid = ?", req.Sid).Select("uf.name, b.ext, b.size, b.url, s.type, s.expired, s.browse").
		Join("LEFT", []string{new(models.UserFileBasic).TableName(), "uf"}, "s.fid = uf.fid").
		Join("LEFT", []string{new(models.BucketBasic).TableName(), "b"}, "uf.rid = b.rid").Get(&data)
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if !has {
		resp.Code = define.SHARE_NOT_EXIST
		return
	}

	resp.Data = data
	return
}
