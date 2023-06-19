package stateful

import (
	docker_hub "github.com/avhaliullin/yandex-alice-k8s-skill/app/docker-hub"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"go.uber.org/zap"
)

type Deps interface {
	GetK8sService() k8s.Service
	GetDockerService() docker_hub.Service
	GetLogger() *zap.Logger
}
