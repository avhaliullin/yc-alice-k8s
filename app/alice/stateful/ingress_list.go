package stateful

import (
	"context"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
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
		return respondText("В каком неймспейсе найти ингрессы?"), nil
	}
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return respondTextF("Я не нашла неймспейс \"%s\"", namespaceName), nil
	}
	ingresses, err := h.k8sService.ListIngresses(ctx, &k8s.ListIngressesReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	if len(ingresses) == 0 {
		return respondTextF("В неймспейсе %s нет ингрессов", namespace), nil
	}
	servicesStr := strings.Join(ingresses, "\n")
	//TODO(plurals)
	return respondTextF("В неймспейсе \"%s\" %d ингрессов:\n%s", namespace, len(ingresses), servicesStr), nil
}
