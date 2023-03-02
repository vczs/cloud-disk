package logic

import (
	"cloud-disk/disk/define"
	"cloud-disk/disk/models"
	"errors"

	"xorm.io/xorm"
)

func FolderSizeChange(session *xorm.Session, action, fid string, size int64, number int) error {
	var pid string
	tempFid := fid
	for pid != "0" {
		res, _ := session.Exec("UPDATE "+new(models.UserFileBasic).TableName()+" SET size = size "+action+" ?,number = number "+action+" ? WHERE fid = ? and type = ?", size, number, tempFid, define.FileTypeFolder)
		affect, err := res.RowsAffected()
		if err != nil || affect < 1 {
			return errors.New("size increase failed" + err.Error())
		}

		has, err := session.Table(new(models.UserFileBasic).TableName()).Where("fid = ? and type = ?", tempFid, define.FileTypeFolder).Select("pid").Get(&pid)
		if err != nil || !has {
			return errors.New("size increase failed" + err.Error())
		}
		tempFid = pid
	}
	return nil
}

func DeleteFile(session *xorm.Session, fileType int, fid, uid string) error {
	if fileType == define.FileTypeFolder {
		// 删除所有子文件
		_, err := session.Where("uid = ? and pid = ? and type = ?", uid, fid, define.FileTypeFile).Delete(new(models.UserFileBasic))
		if err != nil {
			return errors.New("delete file failed")
		}
		// 获取所有子文件夹fid
		var fids []string
		err = session.Table(new(models.UserFileBasic).TableName()).Where("uid = ? and pid = ? and type = ?", uid, fid, define.FileTypeFolder).Select("fid").Find(&fids)
		if err != nil {
			return errors.New("delete file failed")
		}
		// 递归删除获取到的所有子文件夹
		for _, v := range fids {
			err = DeleteFile(session, define.FileTypeFolder, v, uid)
			if err != nil {
				return errors.New("delete file failed")
			}
		}
	}
	// 最后删除目标文件
	affect, err := session.Where("uid = ? and fid = ?", uid, fid).Delete(new(models.UserFileBasic))
	if affect < 1 || err != nil {
		return errors.New("delete file failed")
	}
	return nil
}
