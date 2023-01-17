package config

const (
	SUCCESS = 0 // success

	REQUEST_OFTEN = 10001 // 请求频繁
	SERVER_PANIC  = 10002 // 服务器异常

	USER_NOT_EXIST      = 20001 // 用户不存在
	USER_PASSWORD_ERR   = 20002 // 用户名或密码错误
	USER_NAME_HAS       = 20003 // 用户名已存在
	USER_REGISTER_ERR   = 20004 // 用户注册失败
	USER_PASSWORD_EMPTY = 20005 // 用户名或密码为空
	EMAIL_EMPTY         = 20006 // 邮箱为空
	EMAIL_CODE_EMPTY    = 20007 // 验证码为空

	EMAIL_REGISTERED = 30001 // 该邮箱已被注册
	EMAIL_CODE_ERR   = 30002 // 邮箱验证码错误

	CAP_OVERFLOW         = 40001 // 容量超出上限
	FILE_NOT_EXIST       = 40003 // 文件不存在
	ACCESS_DENIED        = 40004 // 拒绝访问
	FILE_NAME_HAS        = 40005 // 文件名已存在
	FOLDER_NOT_EXIST     = 40006 // 文件夹不存在
	FILE_NAME_EMPTY      = 40007 // 文件名不能为空
	FOLDER_NAME_EMPTY    = 40008 // 文件夹名不能为空
	FILE_RID_EMPTY       = 40009 // 用户文件ID不能为空
	FILE_PID_EMPTY       = 40010 // 目标文件夹ID不能为空
	TARGET_NOT_FOLDER    = 40011 // 目标不是文件夹
	FILE_MOVE_FAILED     = 40012 // 文件移动失败
	FILE_DELETE_FAILED   = 40013 // 文件删除失败
	TARGET_PATH_NAME_HAS = 40014 // 目标目录下文件名已存在
	FOLDER_NUMBER_MAX    = 40015 // 文件夹数量上限

	SHARE_FAILED     = 50001 // 分享失败
	FOLDER_NOT_SHARE = 50002 // 非文件不可分享
	FILE_SID_EMPTY   = 50003 // 分享ID不能为空
	SHARE_NOT_EXIST  = 50004 // 资源不存在
	SHARE_EMPTY      = 50005 // 请输入资源密码
	SHARE_PWD_ERR    = 50006 // 资源密码错误
	SHARE_SAVE_ERR   = 50008 // 资源保存失败
	SHARE_EXPIRE     = 50009 // 资源过期失效

	FILE_MD5_EMPTY       = 60001 // 文件md5不能为空
	FILE_PATCH_NOT_EXIST = 60002 // 文件碎片不存在
	PATCH_UPLOAD_FAILED  = 60003 // 分块上传失败
	COMPLETE_PART_FAILED = 60004 // 完成分块上传失败
	INIT_PART_FAILED     = 60005 // 分块上传初始化失败
	FILE_UPLOAD_FAILED   = 60006 // 文件上传失败
)

var message = map[int]string{
	SUCCESS: "success",

	REQUEST_OFTEN: "请求频繁,请稍后重试!",
	SERVER_PANIC:  "服务器异常!",

	USER_NOT_EXIST:      "用户不存在!",
	USER_PASSWORD_ERR:   "用户名或密码错误!",
	USER_NAME_HAS:       "用户名已存在!",
	USER_REGISTER_ERR:   "用户注册失败!",
	USER_PASSWORD_EMPTY: "用户名或密码为空!",
	EMAIL_EMPTY:         "邮箱为空!",
	EMAIL_CODE_EMPTY:    "验证码为空!",

	EMAIL_REGISTERED: "该邮箱已被注册!",
	EMAIL_CODE_ERR:   "邮箱验证码错误!",

	CAP_OVERFLOW:         "容量超出上限!",
	FILE_NOT_EXIST:       "文件不存在!",
	ACCESS_DENIED:        "权限不足,拒绝访问!",
	FILE_NAME_HAS:        "文件名已存在!",
	FOLDER_NOT_EXIST:     "文件夹不存在!",
	FILE_NAME_EMPTY:      "文件名不能为空!",
	FOLDER_NAME_EMPTY:    "文件夹名不能为空!",
	FILE_RID_EMPTY:       "用户文件ID不能为空!",
	FILE_PID_EMPTY:       "目标文件夹ID不能为空!",
	TARGET_NOT_FOLDER:    "目标不是文件夹!",
	FILE_MOVE_FAILED:     "文件移动失败!",
	FILE_DELETE_FAILED:   "文件删除失败!",
	TARGET_PATH_NAME_HAS: "目标目录下文件名已存在!",
	FOLDER_NUMBER_MAX:    "文件夹数量上限!",

	SHARE_FAILED:     "文件分享失败!",
	FOLDER_NOT_SHARE: "非文件不可分享!",
	FILE_SID_EMPTY:   "分享ID不能为空!",
	SHARE_NOT_EXIST:  "资源不存在!",
	SHARE_EMPTY:      "请输入资源密码!",
	SHARE_PWD_ERR:    "资源密码错误!",
	SHARE_SAVE_ERR:   "资源保存失败!",
	SHARE_EXPIRE:     "资源已过期失效!",

	FILE_MD5_EMPTY:       "文件md5不能为空!",
	FILE_PATCH_NOT_EXIST: "文件碎片不存在!",
	PATCH_UPLOAD_FAILED:  "分块上传失败!",
	COMPLETE_PART_FAILED: "完成分块上传失败!",
	INIT_PART_FAILED:     "分块上传初始化失败!",
	FILE_UPLOAD_FAILED:   "文件上传失败!",
}

// GetMessage 获取message
func GetMessage(code int) string {
	if msg, ok := message[code]; ok {
		return msg
	} else {
		return "服务器发生未知错误~"
	}
}
