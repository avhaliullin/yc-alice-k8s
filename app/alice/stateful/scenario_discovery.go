package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

func (h *Handler) discoverScenarios(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.DiscoverScenarios
	if intnt == nil {
		return nil, nil
	}
	//TODO(full text)
	return respondText("Я умею заглядывать в неймспейсы: считать и искать сломанные поды, могу перечислить сервисы и ингрессы в неймспейсе"), nil
}
