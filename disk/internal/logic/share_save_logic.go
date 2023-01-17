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

type ShareSaveModel struct {
	Name    string `json:"name"`
	Rid     string `json:"rid"`
	Size    int64  `json:"size"`
	Type    int    `json:"type"`
	Encrypt string `json:"encrypt"`
	Expired int64  `json:"expired"`
}

type ShareSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareSaveLogic {
	return &ShareSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareSaveLogic) ShareSave(req *types.ShareSaveRequest, uid string) (resp *types.ShareSaveReply, err error) {
	resp = &types.ShareSaveReply{Code: 0}

	// 参数校验
	if req.Sid == "" {
		resp.Code = config.FILE_SID_EMPTY
		return
	}

	// 检查要保存的分享资源是否存在
	shareSaveModel := new(ShareSaveModel)
	has, err := l.svcCtx.Engine.Table(new(models.ShareBasic).TableName()).Alias("s").Where("sid = ?", req.Sid).
		Select("uf.name, uf.rid, uf.size, s.type, s.encrypt, s.expired").Join("LEFT", []string{new(models.UserFileBasic).TableName(), "uf"}, "uf.fid = s.fid").
		Get(shareSaveModel)
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if !has {
		resp.Code = config.SHARE_NOT_EXIST
		return
	}

	// 检查资源是否在有效期
	if shareSaveModel.Expired < time.Now().Unix() {
		resp.Code = config.SHARE_EXPIRE
		return
	}

	// 检查加密资源密码
	if shareSaveModel.Type == config.ShareTypeEncrypt {
		if req.Encrypt == "" {
			resp.Code = config.SHARE_EMPTY
			return
		} else {
			if req.Encrypt != shareSaveModel.Encrypt {
				resp.Code = config.SHARE_PWD_ERR
				return
			}
		}
	}

	// 检查目标文件夹是否存在  不传表示保存在一级目录
	pid := req.Pid
	if pid == "" || pid == "0" {
		pid = "0"
	} else {
		uf := new(models.UserFileBasic)
		has, err = l.svcCtx.Engine.Where("fid = ? and uid = ?", pid, uid).Select("type").Get(uf)
		if err != nil {
			helper.VczsLog("", err)
			return nil, err
		}
		if !has {
			resp.Code = config.FOLDER_NOT_EXIST
			return
		}
		// 检查目标文件是否为文件夹
		if uf.Type != config.FileTypeFolder {
			resp.Code = config.TARGET_NOT_FOLDER
			return
		}
	}

	// 检查目标文件夹下是否有与要操作的文件 重名的文件
	num, err := l.svcCtx.Engine.Where("pid = ? and name = ? and uid = ?", pid, shareSaveModel.Name, uid).Count(new(models.UserFileBasic))
	if err != nil {
		helper.VczsLog("", err)
		return nil, err
	}
	if num > 0 {
		resp.Code = config.TARGET_PATH_NAME_HAS
		return
	}

	// 构建数据
	fid := helper.GetUid()
	saveuserFile := models.UserFileBasic{
		Fid:    fid,
		Uid:    uid,
		Pid:    pid,
		Name:   shareSaveModel.Name,
		Rid:    shareSaveModel.Rid,
		Type:   config.FileTypeFile,
		Size:   shareSaveModel.Size,
		Number: 1,
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
			resp.Code = config.SHARE_SAVE_ERR
		} else {
			session.Commit()
		}
	}()
	// 保存文件
	affect, err := session.Insert(&saveuserFile)
	if err != nil || affect < 1 {
		panic(config.SHARE_SAVE_ERR)
	}
	// 更新用户文件数量 大小
	if pid != "0" {
		FolderSizeChange(session, config.FolderSizeChangeIncrease, pid, shareSaveModel.Size, 1)
		if err != nil {
			panic(config.SHARE_SAVE_ERR)
		}
	}
	// 更新用户容量
	_, err = session.Exec("UPDATE "+new(models.UserBasic).TableName()+" SET now_volume = now_volume + ? WHERE uid = ?", shareSaveModel.Size, uid)
	if err != nil {
		panic(config.SHARE_SAVE_ERR)
	}

	resp.Data.Fid = fid
	return
}
