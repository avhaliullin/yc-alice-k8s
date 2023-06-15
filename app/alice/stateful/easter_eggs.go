package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

func (h *Handler) easterEggs(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intents := req.Request.NLU.Intents
	if intents.EasterDBLaunch != nil {
		return respondText(
			"Если у вас возникает такой вопрос - то нет",
		), nil
	} else if intents.EasterHowTo != nil {
		return respondText(
			"В Openshift это работает из коробки",
		), nil
	} else if intents.EasterWhatIsK8s != nil {
		return respondText(
			"Kubernetes - это пять бинарей",
		), nil
	}
	return nil, nil
}
