package stateful

import (
	"context"
	"fmt"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
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
	if len(nss) == 0 {
		return respondText("неймспейсы не найдены"), nil
	}
	text := "Я нашла такие неймспейсы:\n"
	for _, ns := range nss {
		text = text + fmt.Sprintf("%s\n", ns)
	}
	return respondText(text), nil
}
