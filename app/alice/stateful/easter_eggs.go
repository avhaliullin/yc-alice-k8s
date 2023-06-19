package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

func (h *Handler) easterEggs(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intents := req.Request.NLU.Intents
	if intents.EasterDBLaunch != nil {
		return resp.EasterDBLaunch(), nil
	} else if intents.EasterHowTo != nil {
		return resp.EasterHowTo(), nil
	} else if intents.EasterWhatIsK8s != nil {
		return resp.EasterWhatIsK8s(), nil
	}
	return nil, nil
}
