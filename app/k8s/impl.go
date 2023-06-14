// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"context"
	"fmt"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

var _ Service = &service{}

type service struct {
	client *kubernetes.Clientset
}

func NewService(deps Deps) (Service, error) {
	err := rest.RegisterAuthProviderPlugin("yciam", authProviderFactory(deps.GetIAMAuth()))
	if err != nil {
		return nil, fmt.Errorf("failed to register iam auth plugin: %w", err)
	}
	appConf := deps.GetConfig()
	client, err := kubernetes.NewForConfig(&rest.Config{
		Host: appConf.K8sHost,
		AuthProvider: &api.AuthProviderConfig{
			Name: "yciam",
		},
		TLSClientConfig: rest.TLSClientConfig{
			CAData: appConf.K8sCA,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("k8s client init failed: %w", err)
	}
	return &service{
		client: client,
	}, nil
}

func (s *service) ListNamespaces(ctx context.Context, req *ListNamespacesReq) ([]string, errors.Err) {
	resp, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.NewInternal(err)
	}
	//TODO(pagination?)
	res := make([]string, len(resp.Items))[:0]
	for _, ns := range resp.Items {
		res = append(res, ns.Name)
	}
	return res, nil
}

func (s *service) CountPods(ctx context.Context, req *CountPodsReq) (int, errors.Err) {
	resp, err := s.client.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, errors.NewInternal(err)
	}
	//TODO(pagination?)
	return resp.Size(), nil
}
