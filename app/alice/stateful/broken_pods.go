package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func (h *Handler) brokenPods(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.BrokenPods
	if intnt == nil {
		return nil, nil
	}

	namespaceName, ok := intnt.Slots.Namespace.AsString()
	if !ok {
		return respondText("В каком неймспейсе поискать сломанные поды?"), nil
	}
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return respondTextF("Я не нашла неймспейс \"%s\"", namespaceName), nil
	}
	statuses, err := h.k8sService.GetPodStatuses(ctx, &k8s.PodStatusesReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	brokenCnt := statuses.Failed + statuses.Unknown
	if brokenCnt == 0 {
		return respondTextF("В неймспейсе %s нет сломанных подов", namespace), nil
	}
	return respondTextF("В неймспейсе \"%s\" %d сломанных подов", namespace, brokenCnt), nil
}
