package stateful

import (
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"go.uber.org/zap"
)

type Deps interface {
	GetK8sService() k8s.Service
	GetLogger() *zap.Logger
}
