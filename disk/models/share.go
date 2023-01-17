package models

import "database/sql"

type ShareBasic struct {
	Id        int            `xorm:"pk autoincr"`
	Sid       string         `xorm:"varchar(50) comment(分享id)"`
	Uid       string         `xorm:"varchar(50) comment(用户id)"`
	Fid       string         `xorm:"varchar(50) comment(用户文件id)"`
	Type      int            `xorm:"smallint comment(分享类型 0:公开 1:加密)"`
	Encrypt   sql.NullString `xorm:"varchar(50) comment(分享密码 null表示公开分享)"`
	Expired   int64          `xorm:"expired comment(分享失效时间)"`
	Browse    int            `xorm:"browse comment(浏览量)"`
	CreatedAt int64          `xorm:"created ct"`
	UpdatedAt int64          `xorm:"updated ut"`
	DeletedAt int64          `xorm:"deleted dt"`
}

func (ShareBasic) TableName() string {
	return "share"
}
