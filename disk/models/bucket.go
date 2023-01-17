package models

type BucketBasic struct {
	Id        int    `xorm:"pk autoincr"`
	Rid       string `xorm:"unique varchar(100) comment(资源id)"`
	Hash      string `xorm:"varchar(100) comment(资源hash值)"`
	Ext       string `xorm:"varchar(30) comment(资源后缀)"`
	Size      int64  `xorm:"size double comment(资源大小)"`
	Url       string `xorm:"varchar(255) comment(资源url)"`
	CreatedAt int64  `xorm:"created ct"`
	UpdatedAt int64  `xorm:"updated ut"`
	DeletedAt int64  `xorm:"deleted dt"`
}

func (BucketBasic) TableName() string {
	return "bucket"
}
