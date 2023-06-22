package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
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
		return resp.AskNSForBrokenPods().
			WithState(&aliceapi.StateData{State: aliceapi.StateBrokenPodsReqNS}), nil
	}
	return h.doBrokenPods(ctx, namespaceName)
}

func (h *Handler) brokenPodsReqNs(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	return h.doBrokenPods(ctx, req.Request.OriginalUtterance)
}

func (h *Handler) doBrokenPods(ctx context.Context, namespaceName string) (*aliceapi.Response, errors.Err) {
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return resp.NSNotFound(namespaceName), nil
	}
	statuses, err := h.k8sService.GetPodStatuses(ctx, &k8s.PodStatusesReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	brokenCnt := statuses.Failed + statuses.Unknown
	return resp.BrokenPodsInNS(namespace, brokenCnt), nil
}
