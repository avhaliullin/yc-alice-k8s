package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
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
		return resp.AskNSForCountingPods().
			WithState(&aliceapi.StateData{State: aliceapi.StateCountPodsReqNS}), nil
	}
	return h.doCountPods(ctx, namespaceName)
}

func (h *Handler) countPodsReqNs(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	return h.doCountPods(ctx, req.Request.OriginalUtterance)
}

func (h *Handler) doCountPods(ctx context.Context, namespaceName string) (*aliceapi.Response, errors.Err) {
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return resp.NSNotFound(namespaceName), nil
	}
	podsCount, err := h.k8sService.CountPods(ctx, &k8s.CountPodsReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return resp.PodsCountInNS(namespace, podsCount), nil
}
