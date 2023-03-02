package logic

import (
	"context"
	"errors"
	"path"
	"strings"

	"cloud-disk/disk/define"
	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type Etags struct {
	Number int    `json:"number"`
	Etag   string `json:"etag"`
}

type CompletePartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompletePartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompletePartLogic {
	return &CompletePartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompletePartLogic) CompletePart(req *types.CompletePartRequest, uid string) (resp *types.CompletePartReply, err error) {
	resp = &types.CompletePartReply{Code: 0}

	// 参数校验
	uploadId := req.UploadId
	if uploadId == "" {
		return nil, errors.New("upload_id is empty")
	}

	// 检验文件碎片是否存在
	userPatch := new(models.UserPatchBasic)
	has, err := l.svcCtx.Engine.Where("upload_id = ? and uid = ? and number = 0", uploadId, uid).Select("patch_key,etag,size").Get(userPatch)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = define.FILE_PATCH_NOT_EXIST
		return
	}

	// 获取所有碎片eTag
	etags := make([]Etags, 0)
	err = l.svcCtx.Engine.Table(new(models.UserPatchBasic).TableName()).Where("upload_id = ? and uid = ? and number > 0", uploadId, uid).Cols("number", "etag").Find(&etags)
	if err != nil {
		return nil, err
	}

	// 发起完成分块上传
	c := make([]cos.Object, 0)
	for _, v := range etags {
		c = append(c, cos.Object{
			ETag:       v.Etag,
			PartNumber: v.Number,
		})
	}
	err = helper.CosCompletePart(userPatch.PatchKey, uploadId, c)
	if err != nil {
		resp.Code = define.COMPLETE_PART_FAILED
		return
	}

	// 查询分块初始化时记录的文件信息
	userFile := new(models.UserFileBasic)
	has, err = l.svcCtx.Engine.Where("rid = ? and uid = ? and type = ?", uploadId, uid, define.FileTypeUndone).Select("name,fid,pid").Get(userFile)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = define.FILE_PATCH_NOT_EXIST
		return
	}

	// 构建数据
	rid := helper.GetUid()
	bucket := &models.BucketBasic{
		Rid:  rid,
		Hash: userPatch.Etag,
		Ext:  path.Ext(userFile.Name),
		Size: userPatch.Size,
		Url:  define.CosUrl + "/" + userPatch.PatchKey,
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
			resp.Code = define.SERVER_PANIC
		} else {
			session.Commit()
		}
	}()

	affect, err := session.Insert(bucket)
	if err != nil || affect < 1 {
		panic(define.SERVER_PANIC)
	}
	affect, err = session.Where("rid = ?", uploadId).Cols("rid", "name", "type", "size", "number").Update(&models.UserFileBasic{Rid: rid, Name: strings.TrimSuffix(userFile.Name, path.Ext(userFile.Name)), Type: define.FileTypeFile, Number: 1, Size: userPatch.Size})
	if err != nil || affect < 1 {
		panic(define.SERVER_PANIC)
	}
	sqlRes, err := session.Exec("UPDATE "+new(models.UserBasic).TableName()+" SET now_volume = now_volume + ? WHERE uid = ?", userPatch.Size, uid)
	affect, err = sqlRes.RowsAffected()
	if err != nil || affect < 1 {
		panic(define.SERVER_PANIC)
	}
	if userFile.Pid != "0" {
		err = FolderSizeChange(session, define.FolderSizeChangeIncrease, userFile.Pid, userPatch.Size, 1)
		if err != nil {
			panic(define.SERVER_PANIC)
		}
	}
	affect, err = session.Where("upload_id = ? and uid = ?", uploadId, uid).Delete(new(models.UserPatchBasic))
	if err != nil || affect < 1 {
		panic(define.SERVER_PANIC)
	}

	resp.Data.Fid = userFile.Fid
	return resp, nil
}
