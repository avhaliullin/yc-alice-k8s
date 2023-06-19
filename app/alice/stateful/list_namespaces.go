package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func (h *Handler) listNamespaces(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.ListNamespaces
	if intnt == nil {
		return nil, nil
	}
	nss, err := h.k8sService.ListNamespaces(ctx, &k8s.ListNamespacesReq{})
	if err != nil {
		return nil, err
	}
	var nsNames []string
	for _, ns := range nss {
		nsNames = append(nsNames, ns)
	}
	return resp.ListNSs(nsNames), nil
}
