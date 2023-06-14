// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/config"
	iam_auth "github.com/avhaliullin/yandex-alice-k8s-skill/app/iam-auth"
)

type Deps interface {
	GetIAMAuth() iam_auth.Service
	GetConfig() *config.Config
}
