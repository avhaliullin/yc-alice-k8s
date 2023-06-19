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
			return respondText("тогда давайте попробуем заново"), nil
		}
		return nil, nil
	}
	action := deleteDeployActionFromState(&req.State.Session)
	action.confirmed = true
	return h.doDeleteDeploy(ctx, action)
}

func (h *Handler) doDeleteDeploy(ctx context.Context, action *deleteDeployAction) (*aliceapi.Response, errors.Err) {
	if action.name == "" {
		return &aliceapi.Response{
			Response: &aliceapi.Resp{Text: "какой депл+ой удалить?"},
			State:    action.toState(aliceapi.StateDeleteDeployReqName),
		}, nil
	}
	if action.deploymentID == "" {
		deployment, err := h.findDeploymentByName(ctx, k8s.DefaultNS, action.name)
		if err != nil {
			return nil, err
		}
		if deployment == nil {
			return respondTextF("Я не нашла депл+оймент %s", action.name), nil
		}
		action.deploymentID = deployment.Name
	}
	if !action.confirmed {
		//TODO(plurals)
		return &aliceapi.Response{
			Response: respWithTTS(fmt.Sprintf(
				"Удаляю деплой %s. Все верно?",
				action.name,
			)),
			State: action.toState(aliceapi.StateDeleteDeployReqConfirm),
		}, nil
	}
	err := h.k8sService.DeleteDeployment(ctx, &k8s.DeleteDeployReq{Name: action.deploymentID})
	if err != nil {
		log.Error(ctx, "deploy scaling failed", zap.Error(err))
		return respondText("не получилось удалить деплой"), nil
	}
	return respondText("деплой удален"), nil
}
