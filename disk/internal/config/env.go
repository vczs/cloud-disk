package config

import "os"

var JwtKey = "vczs-key" // jwt key

var MailPassword = os.Getenv("163pwd") // 邮箱密码
var CodeLength = 6                     // 验证码长度
var CodeExpire = 300                   // 验证码过期时间（s）

var CosUrl = "https://vcz-1258443713.cos.ap-chengdu.myqcloud.com" // COS地址
var TencentSecretId = os.Getenv("TencentSecretId")                // TencentSecretId
var TencentSecretKey = os.Getenv("TencentSecretKey")              // TencentSecretKey
var CosFolderPath = "disk"

var TokenExpire = 7200        // token过期时间
var RefreshTokenExpire = 9600 // refreshToken过期时间

var DefaultVolume = 1000000 // 用户默认容量

var DefaultPageSize = 10 // 分页查询默认每页条数

var MaxUserFolder = 500 // 用户文件夹最大上限

var FileTypeFolder = 0 // 文件类型: 文件夹
var FileTypeFile = 1   // 文件类型: 文件
var FileTypeUndone = 3 // 文件类型: 未完成

var ShareTypePublic = 0  // 分享类型: 公开
var ShareTypeEncrypt = 1 // 分享类型: 加密

var FolderSizeChangeIncrease = "+" // 文件大小改变：增加
var FolderSizeChangeDecrease = "-" // 文件大小改变：减小
