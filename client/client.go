package client

import (
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/client/mysql/hooks"
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/config"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
	"xorm.io/xorm"
)

func NewQimaoFreeMysqlClient() *xorm.EngineGroup {
	var cfg GroupConfig
	if err := config.Get("mysql", "book").Scan(&cfg); err != nil {
		zap.L().Error("get config err", zap.Error(errors.WithStack(err)))
	}
	masterEngine, err := xorm.NewEngine("mysql", cfg.Master.Dsn)
	if err != nil {
		zap.L().Error("NewEngine", zap.Error(err))
	}
	var slaveEngines []*xorm.Engine
	for _, slave := range cfg.Slaves {
		slaveEngine, err := xorm.NewEngine("mysql", slave.Dsn)
		if err != nil {
			zap.L().Error("NewEngine", zap.Error(err))
		}
		slaveEngines = append(slaveEngines, slaveEngine)
	}

	engineGroup, err := xorm.NewEngineGroup(masterEngine, slaveEngines, xorm.RandomPolicy())
	if err != nil {
		zap.L().Error("NewEngineGroup", zap.Error(err))
	}
	if cfg.IsDebug {
		engineGroup.ShowSQL(true)
	}
	if cfg.MaxIdle > 0 {
		engineGroup.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxOpen > 0 {
		engineGroup.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxLifetime > 0 {
		engineGroup.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime))
	}
	dsnCfg, _ := mysql.ParseDSN(cfg.Master.Dsn)
	dsnCfg.Passwd = ""
	engineGroup.AddHook(hooks.NewTracingHook(dsnCfg.FormatDSN()))
	return engineGroup
}

type GroupConfig struct {
	MaxIdle     int `json:"max_idle" mapstructure:"max_idle"`
	MaxOpen     int `json:"max_open" mapstructure:"max_open"`
	MaxLifetime int `json:"max_lifetime" mapstructure:"max_lifetime"`
	Master      struct {
		Dsn string `json:"dsn" mapstructure:"dsn"`
	} `json:"master" mapstructure:"master"`
	Slaves []struct {
		Dsn string `json:"dsn" mapstructure:"dsn"`
	} `json:"slaves" mapstructure:"slaves"`
	IsDebug bool `json:"is_debug" mapstructure:"is_debug"`
}
