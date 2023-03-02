package logic

import (
	"context"

	"cloud-disk/disk/define"
	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest, authorization string) (resp *types.RefreshTokenReply, err error) {
	resp = &types.RefreshTokenReply{Code: 0}

	// 验证authorization
	user, err := helper.AuthToken(authorization)
	if err != nil {
		return
	}

	// 生成token
	token, err := helper.GenerateToken(user.Id, user.Uid, user.Name, define.TokenExpire)
	if err != nil {
		return nil, err
	}

	// 生成refreshToken
	refreshToken, err := helper.GenerateToken(user.Id, user.Uid, user.Name, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	// 返回token和refreshToken
	resp.Data.Token = token
	resp.Data.RefreshToken = refreshToken
	return
}
