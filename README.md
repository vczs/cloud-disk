# cloud-disk

> go-zero + xorm 实现的一个网盘系统

## 系统模块
+ [x] 用户
    - [x] 邮箱注册
    - [x] 密码登录
    - [x] 用户详情
    - [x] 用户容量
    - [x] 刷新token
+ [x] 资源上传
    - [x] 资源上传
    - [x] 资源秒传
    - [x] 资源分块上传
    - [x] 对接腾讯云COS
    - [ ] 对接阿里云OOS
    - [ ] 对接MinIO
+ [x] 用户文件
    - [x] 用户文件关联资源
    - [x] 创建文件夹
    - [x] 文件列表
    - [x] 修改文件名
    - [x] 文件删除
    - [x] 文件移动
    - [ ] 回收站(已删除文件列表、恢复已删除文件)
+ [x] 文件分享
  - [x] 分享文件(共享、加密)
  - [x] 查看分享文件
  - [x] 保存分享文件



```text
1.安装go-zero
go get -u github.com/zeromicro/go-zero

2.安装goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

3.使用goctl命令创建服务
goctl api new disk

4.运行项目
go run disk.go -f etc/disk-api.yaml

5.使用api文件生成代码
goctl api go -api disk.api -dir . -style go_zero

6.第三方地址
腾讯云COS后台：https://console.cloud.tencent.com/cos/bucket
腾讯云COS文档：https://cloud.tencent.com/document/product/436/31215
```