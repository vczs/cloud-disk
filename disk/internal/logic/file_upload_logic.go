package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"
	"strings"

	"cloud-disk/disk/helper"
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, r *http.Request) (resp *types.FileUploadReply, err error) {
	resp = &types.FileUploadReply{Code: 0}

	// 获取用户id
	uid := r.Header.Get("Uid")

	// 获取文件信息
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return
	}
	ext := path.Ext(fileHeader.Filename)
	size := helper.ByteToKb(fileHeader.Size)

	// 文件名非空判断
	name := strings.TrimSuffix(fileHeader.Filename, ext)
	if name == "" {
		resp.Code = config.FILE_NAME_EMPTY
		return
	}

	// 判断文件夹是否存在
	pid := r.Form.Get("pid")
	if pid == "" || pid == "0" {
		pid = "0" // 默认一级文件夹
	} else {
		userFile := new(models.UserFileBasic)
		has, err := l.svcCtx.Engine.Where("fid = ? and uid = ?", pid, uid).Cols("type").Get(userFile)
		if err != nil {
			return nil, err
		}
		if !has {
			resp.Code = config.FOLDER_NOT_EXIST
			return resp, nil
		} else {
			if userFile.Type != config.FileTypeFolder {
				resp.Code = config.TARGET_NOT_FOLDER
				return resp, nil
			}
		}
	}

	// 判断文件夹下该文件名是否已存在
	total, err := l.svcCtx.Engine.Where("name = ? and uid = ? and pid = ?", name, uid, pid).Count(new(models.UserFileBasic))
	if err != nil {
		return
	}
	if total > 0 {
		resp.Code = config.FILE_NAME_HAS
		return resp, nil
	}

	// 判断用户容量上限
	userInfo := new(models.UserBasic)
	has, err := l.svcCtx.Engine.Where("uid = ?", uid).Get(userInfo)
	if err != nil {
		return
	}
	if !has {
		resp.Code = config.USER_NOT_EXIST
		return
	}
	if helper.ByteToKb(size)+userInfo.NowVolume > userInfo.TotalVolume {
		resp.Code = config.CAP_OVERFLOW
		return
	}

	// 判断文件在云端是否已存在
	b := make([]byte, size)
	_, err = file.Read(b)
	if err != nil {
		return
	}
	hash := fmt.Sprintf("%x", md5.Sum(b))
	bucketUrl := new(models.BucketBasic)
	has, err = l.svcCtx.Engine.Where("hash = ?", hash).Select("url,rid").Get(bucketUrl)
	if err != nil {
		return
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
			resp.Code = config.FILE_UPLOAD_FAILED
		} else {
			session.Commit()
		}
	}()

	// 构造文件信息
	fid := helper.GetUid()
	userFile := &models.UserFileBasic{
		Fid:    fid,
		Uid:    uid,
		Pid:    pid,
		Name:   name,
		Rid:    bucketUrl.Rid,
		Type:   config.FileTypeFile,
		Size:   size,
		Number: 1,
	}

	// 如果资源不存在：上传资源并保存资源信息
	if !has {
		rid := helper.GetUid()
		url, err := helper.CosUploadFile(r)
		if err != nil {
			return nil, err
		}
		bucket := &models.BucketBasic{
			Rid:  rid,
			Hash: hash,
			Ext:  ext,
			Size: size,
			Url:  url,
		}
		bucketAffect, err := session.Insert(bucket)
		if err != nil || bucketAffect < 1 {
			panic(config.FILE_UPLOAD_FAILED)
		}
		userFile.Rid = rid
	}

	userFileAffect, err := session.Insert(userFile)
	if err != nil || userFileAffect < 1 {
		panic(config.FILE_UPLOAD_FAILED)
	}

	_, err = session.Exec("UPDATE "+new(models.UserBasic).TableName()+" SET now_volume = now_volume + ? WHERE uid = ?", size, uid)
	if pid != "0" {
		err = FolderSizeChange(session, config.FolderSizeChangeIncrease, pid, size, 1)
		if err != nil {
			panic(config.FILE_UPLOAD_FAILED)
		}
	}

	resp.Data.Fid = fid
	return
}
