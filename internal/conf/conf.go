package conf

import (
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/config"
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/config/encoder/yaml"
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/config/source"
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/config/source/file"
	"codeup.aliyun.com/qimao/go-contrib/prototype/pkg/env"
	"go.uber.org/zap"
)

const PREFIXPATH = "config/"

func InitConfig() {
	runMode := env.Mode()
	err := config.Load(
		file.NewSource(
			file.WithPath(PREFIXPATH+runMode+"/mysql.yaml"),
			source.WithEncoder(yaml.NewEncoder()),
		),
	)
	if err != nil {
		zap.L().Error("load config file fail", zap.Error(err))
	}
}
