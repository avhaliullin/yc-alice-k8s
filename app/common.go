package main

import (
	"context"
	"fmt"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/log"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/util"
	"go.uber.org/zap"
)

func initLogging() (context.Context, error) {
	instanceID := util.GenerateID()

	ctx := context.Background()
	zapConf := zap.NewProductionConfig()
	zapConf.Level.Enabled(zap.DebugLevel)
	zapConf.OutputPaths = []string{"stderr"}
	logger, err := zapConf.Build(zap.AddCallerSkip(3))
	if err != nil {
		return nil, err
	}
	logger = logger.With(zap.String("instanceID", instanceID))
	return log.CtxWithLogger(ctx, logger), nil
}

func assertInitialized(component interface{}, name string) {
	if component == nil {
		panic(fmt.Sprintf("%s wasn't initialized before usage", name))
	}
}
