package logic

import (
	"context"

	"cloud-disk/disk/define"
	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	resp = &types.UserRegisterReply{Code: 0}

	// 参数校验
	if req.Name == "" || req.Password == "" {
		resp.Code = define.USER_PASSWORD_EMPTY
		return
	}
	if req.Email == "" {
		resp.Code = define.EMAIL_EMPTY
		return
	}
	if req.Code == "" {
		resp.Code = define.EMAIL_CODE_EMPTY
		return
	}

	// 判断验证码是否一致
	code, err := l.svcCtx.Redis.Get(l.ctx, req.Email).Result()
	if err != nil || code != req.Code {
		resp.Code = define.EMAIL_CODE_ERR
		return resp, nil
	}
	// 判断用户名是否存在
	num, err := l.svcCtx.Engine.Where("name = ?", req.Name).Count(new(models.UserBasic))
	if err != nil || num > 0 {
		if num > 0 {
			resp.Code = define.USER_NAME_HAS
			return resp, nil
		}
		return nil, err
	}
	// 数据持久化
	user := &models.UserBasic{Uid: helper.GetUid(), Name: req.Name, Password: helper.Md5(req.Password), Email: req.Email, TotalVolume: int64(define.DefaultVolume)}
	num, err = l.svcCtx.Engine.Insert(user)
	if err != nil || num < 1 {
		if num < 1 {
			resp.Code = define.USER_REGISTER_ERR
			return resp, nil
		}
		return nil, err
	}
	return
}
