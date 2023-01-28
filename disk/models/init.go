package models

import (
	"cloud-disk/disk/helper"
	"fmt"

	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// 初始化mysql
func InitMysql(account string, password string, host string, port int, dbname string, charset string) *xorm.Engine {
	// 连接数据库
	connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", account, password, host, port, dbname, charset)
	engine, err := xorm.NewEngine("mysql", connstr)
	if err != nil {
		helper.VczsLog("xorm engine err", err)
		return nil
	}
	// 创建数据库表
	err = engine.Sync(new(UserBasic), new(BucketBasic), new(UserFileBasic), new(ShareBasic), new(UserPatchBasic))
	if err != nil {
		helper.VczsLog("xorm create table err", err)
		return nil
	}
	// // SQL日志
	// xormLogFile, err := os.OpenFile("logs/xorm_sql.log", os.O_APPEND|os.O_WRONLY, 6)
	// if err != nil {
	// 	helper.VczsLog("open xorm_sql.log failed", err)
	// 	return nil
	// }
	// engine.SetLogger(log.NewSimpleLogger(xormLogFile))
	// engine.ShowSQL(true)
	// engine.Logger().SetLevel(log.LOG_INFO)
	return engine
}

// 初始化redis
func InitRedis(host string, port int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
