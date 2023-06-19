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

type scaleDeployAction struct {
	scale        int
	name         string
	deploymentID string
	confirmed    bool
}

func scaleDeployActionFromState(state *aliceapi.StateData) *scaleDeployAction {
	return &scaleDeployAction{
		scale:        state.Scale,
		name:         state.DeployName,
		deploymentID: state.DeployID,
	}
}

func (da *scaleDeployAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:      state,
		Scale:      da.scale,
		DeployName: da.name,
		DeployID:   da.deploymentID,
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
		return resp.ExpectedNumber(), nil
	}
	if scale <= 0 {
		return resp.DeployScaleMinAssert(), nil
	}
	if scale > maxScale {
		return resp.DeployScaleMaxAssert(maxScale), nil
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
			return resp.RejectOnWizard(), nil
		}
		return nil, nil
	}
	action := scaleDeployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doScaleDeploy(ctx, action)
}

func (h *Handler) doScaleDeploy(ctx context.Context, action *scaleDeployAction) (*aliceapi.Response, errors.Err) {
	if action.name == "" {
		return resp.WhichDeployToScale().
			WithState(action.toState(aliceapi.StateScaleDeployReqName)), nil
	}
	if action.deploymentID == "" {
		deployment, err := h.findDeploymentByName(ctx, k8s.DefaultNS, action.name)
		if err != nil {
			return nil, err
		}
		if deployment == nil {
			return resp.DeployNotFound(action.name), nil
		}
		action.deploymentID = deployment.Name
	}
	if !action.confirmed {
		return resp.DeployScalingConfirm(action.name, action.scale).
			WithState(action.toState(aliceapi.StateScaleDeployReqConfirm)), nil
	}
	err := h.k8sService.ScaleDeployment(ctx, &k8s.ScaleDeployReq{
		Name:  action.deploymentID,
		Scale: action.scale,
	})
	if err != nil {
		log.Error(ctx, "deploy scaling failed", zap.Error(err))
		return resp.DeployScalingFail(action.deploymentID), nil
	}
	return resp.DeployScalingSuccess(k8s.DefaultNS, action.deploymentID), nil
}
