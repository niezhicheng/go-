package initalize

import "go.uber.org/zap"

func InitLogeer()  {
	logger,_ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
