package logic

import (
	"context"
	"path"
	"strings"

	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitPartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitPartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitPartLogic {
	return &InitPartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitPartLogic) InitPart(req *types.InitPartRequest, uid string) (resp *types.InitPartReply, err error) {
	resp = &types.InitPartReply{Code: 0}

	// 参数校验
	if req.Md5 == "" {
		resp.Code = config.FILE_MD5_EMPTY
		return
	}
	ext := path.Ext(req.Name)
	name := strings.TrimSuffix(req.Name, ext)
	if name == "" {
		resp.Code = config.FILE_NAME_EMPTY
		return
	}

	// 创建事务
	session := l.svcCtx.Engine.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		err := recover()
		if err != nil {
			session.Rollback()
			resp.Code = config.INIT_PART_FAILED
		} else {
			session.Commit()
		}
	}()

	pid := req.Pid
	if pid == "" {
		pid = "0"
	}
	// 检查目标文件夹是否存在
	if pid != "0" {
		uf := new(models.UserFileBasic)
		has, err := session.Where("fid = ? and uid = ?", pid, uid).Cols("type").Get(uf)
		if err != nil {
			panic(config.INIT_PART_FAILED)
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
	}
	// 检查目标文件夹下是否有与要操作的文件 重名的文件
	num, err := session.Where("pid = ? and name = ? and uid = ?", pid, name, uid).Count(new(models.UserFileBasic))
	if err != nil {
		panic(config.INIT_PART_FAILED)
	}
	if num > 0 {
		resp.Code = config.TARGET_PATH_NAME_HAS
		return resp, nil
	}

	// 构造数据
	fid := helper.GetUid()
	userFile := &models.UserFileBasic{
		Fid:    fid,
		Uid:    uid,
		Pid:    pid,
		Name:   name + ext,
		Type:   config.FileTypeFile,
		Number: 0,
	}

	// 查看文件在云端是否已存在
	bucket := new(models.BucketBasic)
	has, err := session.Where("hash = ?", req.Md5).Select("rid,size").Get(bucket)
	if err != nil {
		panic(config.INIT_PART_FAILED)
	}
	if !has {
		userPatch := new(models.UserPatchBasic)
		uploadIdHas, err := session.Where("etag = ?", req.Md5).Select("upload_id").Get(userPatch)
		if err != nil {
			panic(config.INIT_PART_FAILED)
		}
		if uploadIdHas {
			resp.Data.UploadId = userPatch.UploadId
			return resp, nil
		}
		key, uploadId, err := helper.CosInitPart(ext)
		if err != nil {
			return nil, err
		}
		patch := &models.UserPatchBasic{
			PatchIdentity: helper.GetUid(),
			Uid:           uid,
			PatchKey:      key,
			UploadId:      uploadId,
			Number:        0,
			Etag:          req.Md5,
		}
		affect, err := session.Insert(patch)
		if err != nil {
			panic(config.INIT_PART_FAILED)
		}
		if affect < 1 {
			resp.Code = config.INIT_PART_FAILED
			return resp, nil
		}
		userFile.Rid = uploadId
		userFile.Type = config.FileTypeUndone
		resp.Data.UploadId = uploadId
	} else {
		userFile.Rid = bucket.Rid
		userFile.Size = bucket.Size
		userFile.Number = 1
		session.Exec("UPDATE "+new(models.UserBasic).TableName()+" SET now_volume = now_volume + ? WHERE uid = ?", bucket.Size, uid)
		err := FolderSizeChange(session, config.FolderSizeChangeIncrease, pid, bucket.Size, 1)
		if err != nil {
			panic(config.INIT_PART_FAILED)
		}
		resp.Data.Fid = fid
	}

	affect, err := session.Insert(userFile)
	if err != nil || affect < 1 {
		panic(config.INIT_PART_FAILED)
	}

	return resp, nil
}
