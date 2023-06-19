// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package docker_hub

import (
	"context"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/config"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

type Deps interface {
	GetConfig() *config.Config
}
type Service interface {
	ListImages(ctx context.Context) ([]string, errors.Err)
}

var _ Service = &service{}

type service struct {
	images []string
}

func (s *service) ListImages(ctx context.Context) ([]string, errors.Err) {
	return s.images, nil
}

func NewService(deps Deps) (Service, error) {
	return &service{images: deps.GetConfig().DockerImages}, nil
}
