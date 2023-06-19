package stateful

import (
	"context"
	"fmt"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/log"
	"go.uber.org/zap"
)

type scaleDeployAction struct {
	scale     int
	name      string
	confirmed bool
}

func scaleDeployActionFromState(state *aliceapi.StateData) *scaleDeployAction {
	return &scaleDeployAction{
		scale: state.Scale,
		name:  state.Name,
	}
}

func (da *scaleDeployAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State: state,
		Scale: da.scale,
		Name:  da.name,
	}
}

func (h *Handler) scaleDeployFromScratch(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.ScaleDeploy
	if intnt == nil {
		return nil, nil
	}

	var action scaleDeployAction
	scale, ok := intnt.Slots.Scale.AsInt()
	if ok {
		action.scale = scale
		if scale <= 0 {
			return respondText("не могу задепл+оить меньше одной реплики"), nil
		}
		if scale > maxScale {
			//TODO(plurals)
			return respondTextF("я депл+ою не больше %d реплик", maxScale), nil
		}
	}
	name, ok := intnt.Slots.Name.AsString()
	if ok {
		action.name = name
	}
	return h.doScaleDeploy(ctx, &action)
}

func (h *Handler) scaleDeployReqName(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := scaleDeployActionFromState(&req.State.Session)
	action.name = req.Request.OriginalUtterance
	return h.doScaleDeploy(ctx, action)
}

func (h *Handler) scaleDeployReqScale(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := scaleDeployActionFromState(&req.State.Session)
	ok := false
	var scale int
	for _, entity := range req.Request.NLU.Entities {
		scale, ok = entity.AsInt()
		if ok {
			break
		}
	}
	if !ok {
		return respondText("я вас не поняла, давайте попробуем еще раз"), nil
	}
	if scale <= 0 {
		return respondText("не могу задепл+оить меньше одной реплики"), nil
	}
	if scale > maxScale {
		//TODO(plurals)
		return respondTextF("я депл+ою не больше %d реплик", maxScale), nil
	}
	action.scale = scale
	return h.doScaleDeploy(ctx, action)
}

func (h *Handler) scaleDeployReqConfirm(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	if req.Request.NLU.Intents.Confirm == nil {
		if req.Request.NLU.Intents.Reject != nil {
			return respondText("тогда давайте попробуем заново"), nil
		}
		return nil, nil
	}
	action := scaleDeployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doScaleDeploy(ctx, action)
}

func (h *Handler) doScaleDeploy(ctx context.Context, action *scaleDeployAction) (*aliceapi.Response, errors.Err) {
	if action.name == "" {
		return &aliceapi.Response{
			Response: &aliceapi.Resp{Text: "какой депл+ой отскейлить?"},
			State:    action.toState(aliceapi.StateScaleDeployReqName),
		}, nil
	}
	if !action.confirmed {
		//TODO(plurals)
		return &aliceapi.Response{
			Response: respWithTTS(fmt.Sprintf(
				"Масштабирую деплой %s до %d реплик. Все верно?",
				action.name, action.scale,
			)),
			State: action.toState(aliceapi.StateScaleDeployReqConfirm),
		}, nil
	}
	err := h.k8sService.ScaleDeployment(ctx, &k8s.ScaleDeployReq{
		Name:  ru2id(action.name),
		Scale: action.scale,
	})
	if err != nil {
		log.Error(ctx, "deploy scaling failed", zap.Error(err))
		return respondText("не получилось отмасштабировать деплой"), nil
	}
	return respondText("запустила масштабирование"), nil
}
