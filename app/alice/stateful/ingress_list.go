package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func (h *Handler) listIngresses(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.IngressList
	if intnt == nil {
		return nil, nil
	}

	namespaceName, ok := intnt.Slots.Namespace.AsString()
	if !ok {
		return resp.AskNSForIngresses(), nil
	}
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return resp.NSNotFound(namespaceName), nil
	}
	ingresses, err := h.k8sService.ListIngresses(ctx, &k8s.ListIngressesReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return resp.ListIngresses(namespace, ingresses), nil
}
