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

const maxScale = 10

type deployAction struct {
	imageText string
	image     string
	scale     int
	name      string
	confirmed bool
}

func deployActionFromState(state *aliceapi.StateData) *deployAction {
	return &deployAction{
		imageText: state.Image,
		image:     state.ImageID,
		scale:     state.Scale,
		name:      state.DeployName,
	}
}

func (da *deployAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:      state,
		Image:      da.imageText,
		ImageID:    da.image,
		Scale:      da.scale,
		DeployName: da.name,
	}
}

func (h *Handler) deployFromScratch(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.Deploy
	if intnt == nil {
		return nil, nil
	}

	action := deployAction{scale: 1}
	imageText, ok := intnt.Slots.Image.AsString()
	if ok {
		action.imageText = imageText
	}
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
	return h.doDeploy(ctx, &action)
}

func (h *Handler) deployReqImage(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := deployActionFromState(&req.State.Session)
	action.image = ""
	action.imageText = req.Request.OriginalUtterance
	return h.doDeploy(ctx, action)
}

func (h *Handler) deployReqName(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := deployActionFromState(&req.State.Session)
	action.name = req.Request.OriginalUtterance
	return h.doDeploy(ctx, action)
}

func (h *Handler) deployReqConfirm(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	if req.Request.NLU.Intents.Confirm == nil {
		if req.Request.NLU.Intents.Reject != nil {
			return respondText("тогда давайте попробуем заново"), nil
		}
		return nil, nil
	}
	action := deployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doDeploy(ctx, action)
}

func (h *Handler) doDeploy(ctx context.Context, action *deployAction) (*aliceapi.Response, errors.Err) {
	if action.imageText == "" {
		return &aliceapi.Response{
			Response: respWithTTS("какой образ задепл+оить?"),
			State:    action.toState(aliceapi.StateDeployReqImage),
		}, nil
	}
	if action.image == "" {
		image, err := h.findImageByName(ctx, action.imageText)
		if err != nil {
			return nil, err
		}
		if image == "" {
			return respondTextF("я не знаю образ %s", action.imageText), nil
		}
		action.image = image
		action.imageText = image
	}
	if action.name == "" {
		return &aliceapi.Response{
			Response: &aliceapi.Resp{Text: "Как назовем деплой?"},
			State:    action.toState(aliceapi.StateDeployReqName),
		}, nil
	}
	if !action.confirmed {
		//TODO(plurals)
		return &aliceapi.Response{
			Response: respWithTTS(fmt.Sprintf(
				"Запускаю деплой %s из образа %s на %d реплик. Все верно?",
				action.name, action.image, action.scale,
			)),
			State: action.toState(aliceapi.StateDeployConfirm),
		}, nil
	}
	err := h.k8sService.Deploy(ctx, &k8s.DeployReq{
		Image: action.image,
		Name:  ru2id(action.name),
		Scale: action.scale,
	})
	if err != nil {
		log.Error(ctx, "deploy failed", zap.Error(err))
		return respondText("не получилось запустить деплой"), nil
	}
	return respondText("запустила"), nil
}
