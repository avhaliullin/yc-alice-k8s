package main

import (
	"context"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice"
	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/stateful"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/cloud"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/config"
	docker_hub "github.com/avhaliullin/yandex-alice-k8s-skill/app/docker-hub"
	iam_auth "github.com/avhaliullin/yandex-alice-k8s-skill/app/iam-auth"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/log"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/zap"
)

type aliceApp struct {
	ctx           context.Context
	logger        *zap.Logger
	conf          *config.Config
	sdk           *ycsdk.SDK
	handler       alice.Handler
	iamAuth       iam_auth.Service
	k8sService    k8s.Service
	dockerService docker_hub.Service
}

func (a *aliceApp) GetLogger() *zap.Logger {
	assertInitialized(a.logger, "logger")
	return a.logger
}

func (a *aliceApp) GetContext() context.Context {
	assertInitialized(a.ctx, "ctx")
	return a.ctx
}

func (a *aliceApp) GetCloudSDK() *ycsdk.SDK {
	assertInitialized(a.sdk, "sdk")
	return a.sdk
}

func (a *aliceApp) GetIAMAuth() iam_auth.Service {
	return a.iamAuth
}

func (a *aliceApp) GetK8sService() k8s.Service {
	return a.k8sService
}

func (a *aliceApp) GetDockerService() docker_hub.Service {
	return a.dockerService
}

func (a *aliceApp) GetConfig() *config.Config {
	assertInitialized(a.conf, "conf")
	return a.conf
}

var aliceAppInstance *aliceApp

func initAliceApp() (*aliceApp, error) {
	ctx, err := initLogging()
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "initializing alice app")

	aliceAppInstance = &aliceApp{ctx: ctx, conf: config.LoadFromEnv(), logger: log.FromCtx(ctx)}
	aliceAppInstance.sdk, err = cloud.NewSDK(aliceAppInstance)
	if err != nil {
		return nil, err
	}

	aliceAppInstance.iamAuth, err = iam_auth.NewMetadata()
	if err != nil {
		return nil, err
	}
	aliceAppInstance.dockerService, err = docker_hub.NewService(aliceAppInstance)
	if err != nil {
		return nil, err
	}
	aliceAppInstance.k8sService, err = k8s.NewService(aliceAppInstance)
	if err != nil {
		return nil, err
	}
	aliceAppInstance.handler, err = stateful.NewHandler(aliceAppInstance)
	if err != nil {
		return nil, err
	}
	return aliceAppInstance, nil
}

func getAliceApp() (*aliceApp, error) {
	if aliceAppInstance == nil {
		return initAliceApp()
	}
	return aliceAppInstance, nil
}

func AliceHandler(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, error) {
	aliceApp, err := getAliceApp()
	if err != nil {
		return nil, err
	}
	return aliceApp.handler.Handle(ctx, req)
}
