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

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.UserLoginReply, err error) {
	resp = &types.UserLoginReply{Code: config.SUCCESS}

	// 参数校验
	if req.Name == "" || req.Password == "" {
		resp.Code = config.USER_PASSWORD_EMPTY
		return
	}

	// 查询用户
	user := new(models.UserBasic)
	has, err := l.svcCtx.Engine.Where("name = ? and password = ?", req.Name, helper.Md5(req.Password)).Get(user)
	if !has || err != nil {
		if !has {
			resp.Code = config.USER_PASSWORD_ERR
			return
		}
		return
	}
	// 生成token
	token, err := helper.GenerateToken(user.Id, user.Uid, user.Name, config.TokenExpire)
	if err != nil {
		return
	}
	// 生成refreshToken
	refreshToken, err := helper.GenerateToken(user.Id, user.Uid, user.Name, config.RefreshTokenExpire)
	if err != nil {
		return
	}
	// 返回token和refreshToken
	resp.Data.Token = token
	resp.Data.RefreshToken = refreshToken
	return
}
