package logic

import (
	"context"
	"time"

	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCodeLogic {
	return &EmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailCodeLogic) EmailCode(req *types.EmailCodeRequest) (resp *types.EmailCodeReply, err error) {
	resp = &types.EmailCodeReply{Code: 0}

	// 参数校验
	if req.Email == "" {
		resp.Code = config.EMAIL_EMPTY
		return
	}

	num, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(models.UserBasic))
	if err != nil || num > 0 {
		if num > 0 {
			resp.Code = config.EMAIL_REGISTERED
			return resp, nil
		}
		return nil, err
	}
	down, err := l.svcCtx.Redis.TTL(l.ctx, req.Email).Result()
	if err != nil {
		return nil, err
	}
	if down.Seconds() > 0 || down.Seconds() == -1 {
		resp.Code = config.REQUEST_OFTEN
		return resp, nil
	}
	// 获取验证码
	code := helper.GetEmailCode()
	// 发送验证码
	err = helper.SendEmailCode(req.Email, code)
	if err != nil {
		return nil, err
	}
	// 存储验证码
	l.svcCtx.Redis.Set(l.ctx, req.Email, code, time.Second*time.Duration(config.CodeExpire))
	return
}
