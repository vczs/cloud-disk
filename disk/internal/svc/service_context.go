package svc

import (
	"cloud-disk/disk/internal/config"
	"cloud-disk/disk/internal/middleware"
	"cloud-disk/disk/models"

	"github.com/go-redis/redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config
	Engine *xorm.Engine
	Redis  *redis.Client
	Auth   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Engine: models.InitMysql(c.Mysql.Account, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Name, c.Mysql.Charset),
		Redis:  models.InitRedis(c.Redis.Host, c.Redis.Port),
		Auth:   middleware.NewAuthMiddleware().AuthMiddlewareHandle,
	}
}
