package global

import (
	"goBTC/client"
	"goBTC/models"

	"go.uber.org/zap"
)

var (
	CONFIG    *models.Server
	LOG       *zap.Logger
	Client    *client.BTCClient
	MysqlFlag bool
)
