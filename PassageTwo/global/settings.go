package global

import (
	"PassageTwo/pkg/logger"
	"PassageTwo/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DataBaseSetting *setting.DataBaseSettings

	Logger *logger.Logger
)
