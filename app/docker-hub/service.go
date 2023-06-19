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
	ListImages(ctx context.Context) ([]Image, errors.Err)
}

var _ Service = &service{}

type service struct {
	images []Image
}

func (s *service) ListImages(ctx context.Context) ([]Image, errors.Err) {
	return s.images, nil
}

func NewService(deps Deps) (Service, error) {
	var images []Image
	usedNames := make(map[string]bool)
	for pron, name := range deps.GetConfig().DockerImages {
		images = append(images, Image{
			PronouncedName: pron,
			Name:           name,
		})
		if !usedNames[name] {
			usedNames[name] = true
			images = append(images, Image{
				PronouncedName: name,
				Name:           name,
			})
		}
	}
	return &service{images: images}, nil
}

type Image struct {
	PronouncedName string
	Name           string
}
