package models

type UserPatchBasic struct {
	Id            int    `xorm:"pk autoincr"`
	PatchIdentity string `xorm:"patch_identity unique varchar(50) comment(资源碎片id)"`
	Uid           string `xorm:"varchar(50) comment(用户id)"`
	PatchKey      string `xorm:"patch_key comment(资源碎片key)"`
	UploadId      string `xorm:"upload_id varchar(100) comment(资源碎片上传id)"`
	Number        int    `xorm:"smallint comment(资源碎片上传序号)"`
	Etag          string `xorm:"varchar(100) comment(资源碎片tag)"`
	Size          int64  `xorm:"size double comment(资源碎片大小)"`
	CreatedAt     int64  `xorm:"created ct"`
	UpdatedAt     int64  `xorm:"updated ut"`
	DeletedAt     int64  `xorm:"deleted dt"`
}

func (UserPatchBasic) TableName() string {
	return "user_patch"
}
