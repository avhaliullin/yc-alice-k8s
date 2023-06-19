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

type deleteDeployAction struct {
	name         string
	deploymentID string
	confirmed    bool
}

func deleteDeployActionFromState(state *aliceapi.StateData) *deleteDeployAction {
	return &deleteDeployAction{
		name:         state.DeployName,
		deploymentID: state.DeployID,
	}
}

func (dda *deleteDeployAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:      state,
		DeployName: dda.name,
		DeployID:   dda.deploymentID,
	}
}

func (h *Handler) deleteDeployFromScratch(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.DeleteDeploy
	if intnt == nil {
		return nil, nil
	}

	var action deleteDeployAction
	name, ok := intnt.Slots.Name.AsString()
	if ok {
		action.name = name
	}
	return h.doDeleteDeploy(ctx, &action)
}

func (h *Handler) deleteDeployReqName(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := deleteDeployActionFromState(&req.State.Session)
	action.name = req.Request.OriginalUtterance
	return h.doDeleteDeploy(ctx, action)
}

func (h *Handler) deleteDeployReqConfirm(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	if req.Request.NLU.Intents.Confirm == nil {
		if req.Request.NLU.Intents.Reject != nil {
			return resp.RejectOnWizard(), nil
		}
		return nil, nil
	}
	action := deleteDeployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doDeleteDeploy(ctx, action)
}

func (h *Handler) doDeleteDeploy(ctx context.Context, action *deleteDeployAction) (*aliceapi.Response, errors.Err) {
	if action.name == "" {
		return resp.WhichDeployToDelete().
			WithState(action.toState(aliceapi.StateDeleteDeployReqName)), nil
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
		return resp.ConfirmDeletingDeploy(action.name).
			WithState(action.toState(aliceapi.StateDeleteDeployReqConfirm)), nil
	}
	err := h.k8sService.DeleteDeployment(ctx, &k8s.DeleteDeployReq{Name: action.deploymentID})
	if err != nil {
		log.Error(ctx, "deploy scaling failed", zap.Error(err))
		return resp.DeployDeletionFailed(action.deploymentID), nil
	}
	return resp.DeployDeleted(action.deploymentID), nil
}
