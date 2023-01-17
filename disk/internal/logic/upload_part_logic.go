package logic

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPartLogic {
	return &UploadPartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadPartLogic) UploadPart(req *types.UploadPartRequest, r *http.Request) (resp *types.UploadPartReply, err error) {
	resp = &types.UploadPartReply{Code: 0}

	// 参数必填判断
	uploadId := r.PostForm.Get("upload_id")
	if uploadId == "" {
		return nil, errors.New("upload_id is empty")
	}
	num := r.PostForm.Get("number")
	if num == "" {
		return nil, errors.New("number is empty")
	}
	number, err := strconv.Atoi(num)
	if err != nil {
		return nil, err
	}

	// 获取文件信息
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	// 获取uid
	uid := r.Header.Get("Uid")

	// 检验文件碎片是否存在
	userPatch := new(models.UserPatchBasic)
	has, err := l.svcCtx.Engine.Where("upload_id = ? and uid = ? and number = ?", uploadId, uid, 0).Select("patch_key,size").Get(userPatch)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = config.FILE_PATCH_NOT_EXIST
		return
	}

	// 判断用户容量上限
	userInfo := new(models.UserBasic)
	has, err = l.svcCtx.Engine.Where("uid = ?", uid).Select("now_volume,total_volume").Get(userInfo)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = config.USER_NOT_EXIST
		return
	}
	if helper.ByteToKb(fileHeader.Size)+userPatch.Size+userInfo.NowVolume > userInfo.TotalVolume {
		resp.Code = config.CAP_OVERFLOW
		return
	}

	// 上传文件分块
	etag, err := helper.CosUploadPart(r, userPatch.PatchKey, uploadId, number)
	if err != nil {
		return nil, err
	}

	// 构建数据
	patch := &models.UserPatchBasic{
		PatchIdentity: helper.GetUid(),
		Uid:           uid,
		PatchKey:      userPatch.PatchKey,
		UploadId:      uploadId,
		Number:        number,
		Etag:          etag,
		Size:          helper.ByteToKb(fileHeader.Size),
	}

	// 创建事务
	session := l.svcCtx.Engine.NewSession()
	defer session.Close()
	// 开启事务
	if err = session.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		err := recover()
		if err != nil {
			session.Rollback()
			resp.Code = config.PATCH_UPLOAD_FAILED
		} else {
			session.Commit()
		}
	}()
	affect, err := session.Insert(patch)
	if err != nil || affect < 1 {
		panic(config.FILE_MOVE_FAILED)
	}
	affect, err = session.Where("upload_id = ? and number = 0", uploadId).Cols("size").Update(&models.UserPatchBasic{Size: userPatch.Size + helper.ByteToKb(fileHeader.Size)})
	if err != nil || affect < 1 {
		panic(config.FILE_MOVE_FAILED)
	}

	return resp, nil
}
