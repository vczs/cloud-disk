package logic

import (
	"context"

	"cloud-disk/disk/define"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"
	"cloud-disk/disk/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.FileDeleteRequest, uid string) (resp *types.FileDeleteReply, err error) {
	resp = &types.FileDeleteReply{Code: 0}

	// 参数校验
	fid := req.Fid
	if fid == "" {
		resp.Code = define.FILE_RID_EMPTY
		return
	}

	// 检查文件是否存在
	userFile := new(models.UserFileBasic)
	has, err := l.svcCtx.Engine.Where("uid = ? and fid = ?", uid, fid).Select("pid,size,number,type").Get(userFile)
	if err != nil {
		return nil, err
	}
	if !has {
		resp.Code = define.FILE_NOT_EXIST
		return resp, nil
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
			resp.Code = define.FILE_DELETE_FAILED
		} else {
			session.Commit()
		}
	}()

	// 从数据库中软删除
	if userFile.Pid != "0" {
		err = FolderSizeChange(session, define.FolderSizeChangeDecrease, userFile.Pid, userFile.Size, userFile.Number)
		if err != nil {
			panic(define.FILE_DELETE_FAILED)
		}
	}
	err = DeleteFile(session, userFile.Type, fid, uid)
	if err != nil {
		panic(define.FILE_DELETE_FAILED)
	}
	sqlRes, err := session.Exec("UPDATE "+new(models.UserBasic).TableName()+" SET now_volume = now_volume - ? WHERE uid = ?", userFile.Size, uid)
	if err != nil {
		panic(define.FILE_DELETE_FAILED)
	}
	affect, err := sqlRes.RowsAffected()
	if affect < 1 || err != nil {
		panic(define.FILE_DELETE_FAILED)
	}

	return
}
