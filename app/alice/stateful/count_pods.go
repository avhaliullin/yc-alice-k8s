package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func (h *Handler) countPods(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.CountPods
	if intnt == nil {
		return nil, nil
	}

	namespaceName, ok := intnt.Slots.Namespace.AsString()
	if !ok {
		return respondText("В каком неймспейсе посчитать поды?"), nil
	}
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return respondTextF("Я не нашла неймспейс \"%s\"", namespaceName), nil
	}
	podsCount, err := h.k8sService.CountPods(ctx, &k8s.CountPodsReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return respondTextF("В неймспейсе \"%s\" %d подов", namespace, podsCount), nil
}
