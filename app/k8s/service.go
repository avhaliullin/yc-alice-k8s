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
	GetPodStatuses(ctx context.Context, req *PodStatusesReq) (*PodStatusesResp, errors.Err)
	ListServices(ctx context.Context, req *ListServicesReq) ([]string, errors.Err)
	ListIngresses(ctx context.Context, req *ListIngressesReq) ([]string, errors.Err)
	Deploy(ctx context.Context, req *DeployReq) errors.Err
}

type ListNamespacesReq struct {
}

type CountPodsReq struct {
	Namespace string
}

type PodStatusesReq struct {
	Namespace string
}

type PodStatusesResp struct {
	Pending   int
	Running   int
	Succeeded int
	Failed    int
	Unknown   int
}

type ListServicesReq struct {
	Namespace string
}

type ListIngressesReq struct {
	Namespace string
}

type DeployReq struct {
	Image string
	Name  string
	Scale int
}
