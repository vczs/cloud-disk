package logic

import (
	"context"

	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(uid string) (resp *types.UserInfoReply, err error) {
	resp = &types.UserInfoReply{Code: 0}
	// 查询用户
	user := new(models.UserBasic)
	has, err := l.svcCtx.Engine.Where("uid = ?", uid).Get(user)
	if err != nil {
		return
	}
	if !has {
		resp.Code = config.USER_NOT_EXIST
		return
	}
	// 返回用户详情
	resp.Data = types.UserInfoReplyData{Name: user.Name, Email: user.Email, NowVolume: user.NowVolume, TotalVolume: user.TotalVolume}
	return
}
