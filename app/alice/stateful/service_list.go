package stateful

import (
	"context"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func (h *Handler) listServices(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.ServiceList
	if intnt == nil {
		return nil, nil
	}

	namespaceName, ok := intnt.Slots.Namespace.AsString()
	if !ok {
		return respondText("В каком неймспейсе найти сервисы?"), nil
	}
	namespace, err := h.findNamespaceByName(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return respondTextF("Я не нашла неймспейс \"%s\"", namespaceName), nil
	}
	services, err := h.k8sService.ListServices(ctx, &k8s.ListServicesReq{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return respondTextF("В неймспейсе %s нет сервисов", namespace), nil
	}
	servicesStr := strings.Join(services, "\n")
	return respondTextF("В неймспейсе \"%s\" %d сервисов:\n%s", namespace, len(services), servicesStr), nil
}
