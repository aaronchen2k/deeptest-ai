package _logs

import (
	"go.uber.org/zap"
)

var Zap *zap.Logger

// init in other places for server and agent
func SetLogger(val *zap.Logger) {
	Zap = val
}
