// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"context"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

type Service interface {
	ListNamespaces(ctx context.Context, req *ListNamespacesReq) ([]string, errors.Err)
	CountPods(ctx context.Context, req *CountPodsReq) (int, errors.Err)
}

type ListNamespacesReq struct {
}

type CountPodsReq struct {
	Namespace string
}
