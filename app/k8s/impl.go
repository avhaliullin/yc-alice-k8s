// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"context"
	"fmt"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

const DefaultNS = "default"

var _ Service = &service{}

type service struct {
	client *kubernetes.Clientset
}

func (s *service) DeleteDeployment(ctx context.Context, req *DeleteDeployReq) errors.Err {
	err := s.client.AppsV1().Deployments(DefaultNS).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return errors.NewInternal(err)
	}
	return nil
}

func (s *service) ScaleDeployment(ctx context.Context, req *ScaleDeployReq) errors.Err {
	deployments := s.client.AppsV1().Deployments(DefaultNS)
	scale, err := deployments.GetScale(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return errors.NewInternal(err)
	}
	scale.Spec.Replicas = int32(req.Scale)
	_, err = s.client.AppsV1().Deployments(DefaultNS).UpdateScale(ctx, req.Name, scale, metav1.UpdateOptions{})
	if err != nil {
		return errors.NewInternal(err)
	}
	return nil
}

func (s *service) ListDeployments(ctx context.Context, req *ListDeploymentsReq) ([]appsv1.Deployment, errors.Err) {
	resp, err := s.client.AppsV1().Deployments(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.NewInternal(err)
	}
	//TODO(paging)
	return resp.Items, nil
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
	return len(resp.Items), nil
}

func (s *service) GetPodStatuses(ctx context.Context, req *PodStatusesReq) (*PodStatusesResp, errors.Err) {
	resp, err := s.client.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.NewInternal(err)
	}
	//TODO(pagination)
	var result PodStatusesResp
	for _, pod := range resp.Items {
		switch pod.Status.Phase {
		case corev1.PodPending:
			result.Pending++
		case corev1.PodRunning:
			result.Running++
		case corev1.PodSucceeded:
			result.Succeeded++
		case corev1.PodFailed:
			result.Failed++
		case corev1.PodUnknown:
			result.Unknown++
		}
	}
	return &result, nil
}

func (s *service) ListServices(ctx context.Context, req *ListServicesReq) ([]string, errors.Err) {
	resp, err := s.client.CoreV1().Services(req.Namespace).List(ctx, metav1.ListOptions{})
	//TODO(pagination)
	if err != nil {
		return nil, errors.NewInternal(err)
	}
	result := make([]string, len(resp.Items))[:0]
	for _, service := range resp.Items {
		result = append(result, service.Name)
	}
	return result, nil
}

func (s *service) ListIngresses(ctx context.Context, req *ListIngressesReq) ([]string, errors.Err) {
	resp, err := s.client.NetworkingV1().Ingresses(req.Namespace).List(ctx, metav1.ListOptions{})
	//TODO(pagination)
	if err != nil {
		return nil, errors.NewInternal(err)
	}
	result := make([]string, len(resp.Items))[:0]
	for _, ingress := range resp.Items {
		result = append(result, ingress.Name)
	}
	return result, nil
}

func (s *service) Deploy(ctx context.Context, req *DeployReq) errors.Err {
	_, err := s.client.AppsV1().Deployments(DefaultNS).Create(ctx, &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(req.Scale)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: req.Image, Name: req.Name},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": req.Name},
				},
			},
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app": req.Name,
			}},
		},
		Status: appsv1.DeploymentStatus{},
	}, metav1.CreateOptions{})
	if err != nil {
		return errors.NewInternal(err)
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
