package models

type UserBasic struct {
	Id          int    `xorm:"pk autoincr"`
	Uid         string `xorm:"unique varchar(36) comment(用户id)"`
	Name        string `xorm:"varchar(50) comment(用户名)"`
	Password    string `xorm:"varchar(50) comment(用户密码)"`
	Email       string `xorm:"varchar(100) comment(用户邮箱)"`
	NowVolume   int64  `xorm:"now_volume comment(用户当前容量大小 单位:kb)"`
	TotalVolume int64  `xorm:"total_volume comment(用户云盘总容量大小 单位:kb)"`
	CreatedAt   int64  `xorm:"created ct"`
	UpdatedAt   int64  `xorm:"updated ut"`
	DeletedAt   int64  `xorm:"deleted dt"`
}

func (UserBasic) TableName() string {
	return "user"
}
