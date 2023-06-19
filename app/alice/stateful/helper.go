package stateful

import (
	"context"
	"strings"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/essentialkaos/translit/v2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func (h *Handler) findNamespaceByName(ctx context.Context, nsName string) (string, errors.Err) {
	nss, err := h.k8sService.ListNamespaces(ctx, &k8s.ListNamespacesReq{})
	if err != nil {
		return "", err
	}

	result, ok := findIdByName(nsName, nss)
	if !ok {
		return "", nil
	}
	return result, nil
}

func (h *Handler) findImageByName(ctx context.Context, imageName string) (string, errors.Err) {
	images, err := h.dockerService.ListImages(ctx)
	if err != nil {
		return "", err
	}

	result, ok := text.BestMatch(imageName, text.ImageListMatcher(images))
	if !ok {
		return "", nil
	}
	return images[result].Name, nil
}

func (h *Handler) findDeploymentByName(ctx context.Context, namespace string, deployName string) (*appsv1.Deployment, errors.Err) {
	deployments, err := h.k8sService.ListDeployments(ctx, &k8s.ListDeploymentsReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	result, ok := text.BestMatch(ru2id(deployName), text.DeploymentsMatcher(deployments))
	if !ok {
		return nil, nil
	}
	return &deployments[result], nil
}

func findIdByName(name string, ids []string) (string, bool) {
	idx, ok := text.BestMatch(ru2id(name), text.IDListMatcher(ids))
	if !ok {
		return "", false
	}
	return ids[idx], true
}

func ru2id(text string) string {
	return strings.ReplaceAll(translit.ICAO(strings.ToLower(text)), " ", "-")
}

func mustToJSON(x interface{}) string {
	res, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(res)
}
