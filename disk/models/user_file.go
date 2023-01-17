package models

type UserFileBasic struct {
	Id        int    `xorm:"pk autoincr"`
	Fid       string `xorm:"unique varchar(50) comment(用户文件id)"`
	Uid       string `xorm:"varchar(50) comment(用户id)"`
	Pid       string `xorm:"varchar(50) comment(用户文件父id)"`
	Name      string `xorm:"varchar(100) comment(用户文件名)"`
	Rid       string `xorm:"varchar(100) comment(文件在资源桶id)"`
	Type      int    `xorm:"smallint comment(用户文件类型 0:文件夹 1:文件)"`
	Size      int64  `xorm:"size double comment(用户文件大小)"`
	Number    int    `xorm:"smallint comment(用户当前文件下数量)"`
	CreatedAt int64  `xorm:"created ct"`
	UpdatedAt int64  `xorm:"updated ut"`
	DeletedAt int64  `xorm:"deleted dt"`
}

func (UserFileBasic) TableName() string {
	return "user_file"
}
