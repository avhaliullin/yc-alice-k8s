package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/log"
	"go.uber.org/zap"
)

const maxScale = 10

type deployAction struct {
	imageName string
	imageID   string
	scale     int
	name      string
	confirmed bool
}

func deployActionFromState(state *aliceapi.StateData) *deployAction {
	return &deployAction{
		imageName: state.Image,
		imageID:   state.ImageID,
		scale:     state.Scale,
		name:      state.DeployName,
	}
}

func (da *deployAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:      state,
		Image:      da.imageName,
		ImageID:    da.imageID,
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
		action.imageName = imageText
	}
	scale, ok := intnt.Slots.Scale.AsInt()
	if ok {
		action.scale = scale
		if scale <= 0 {
			return resp.DeployScaleMinAssert(), nil
		}
		if scale > maxScale {
			return resp.DeployScaleMaxAssert(maxScale), nil
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
	action.imageID = ""
	action.imageName = req.Request.OriginalUtterance
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
			return resp.RejectOnWizard(), nil
		}
		return nil, nil
	}
	action := deployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doDeploy(ctx, action)
}

func (h *Handler) doDeploy(ctx context.Context, action *deployAction) (*aliceapi.Response, errors.Err) {
	if action.imageName == "" {
		return resp.WhichImageToDeploy().
			WithState(action.toState(aliceapi.StateDeployReqImage)), nil
	}
	if action.imageID == "" {
		image, err := h.findImageByName(ctx, action.imageName)
		if err != nil {
			return nil, err
		}
		if image == "" {
			return resp.ImageNotFound(action.imageName), nil
		}
		action.imageID = image
		action.imageName = image
	}
	if action.name == "" {
		return resp.HowToNameDeploy().WithState(action.toState(aliceapi.StateDeployReqName)), nil
	}
	if !action.confirmed {
		return resp.ConfirmDeploy(action.name, action.imageID, action.scale).
			WithState(action.toState(aliceapi.StateDeployConfirm)), nil
	}
	err := h.k8sService.Deploy(ctx, &k8s.DeployReq{
		Image: action.imageID,
		Name:  ru2id(action.name),
		Scale: action.scale,
	})
	if err != nil {
		log.Error(ctx, "deploy failed", zap.Error(err))
		return resp.DeployFailed(), nil
	}
	return resp.DeployStarted(), nil
}
