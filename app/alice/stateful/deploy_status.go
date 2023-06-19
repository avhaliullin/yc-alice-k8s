package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

type deployStatusAction struct {
	name        string
	namespace   string
	namespaceID string
}

func deployStatusActionFromState(state *aliceapi.StateData) *deployStatusAction {
	return &deployStatusAction{
		name:        state.Name,
		namespace:   state.Namespace,
		namespaceID: state.NamespaceID,
	}
}

func (dsa *deployStatusAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:       state,
		Name:        dsa.name,
		Namespace:   dsa.namespace,
		NamespaceID: dsa.namespaceID,
	}
}

func (h *Handler) deployStatusFromScratch(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	intnt := req.Request.NLU.Intents.DeployStatus
	if intnt == nil {
		return nil, nil
	}

	var action deployStatusAction
	name, ok := intnt.Slots.Name.AsString()
	if ok {
		action.name = name
	}
	namespace, ok := intnt.Slots.Namespace.AsString()
	if ok {
		action.namespace = namespace
	}
	return h.doDeployStatus(ctx, &action)
}

func (h *Handler) deployStatusReqNamespace(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := deployStatusActionFromState(&req.State.Session)
	action.namespace = req.Request.OriginalUtterance
	return h.doDeployStatus(ctx, action)
}

func (h *Handler) deployStatusReqName(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	action := deployStatusActionFromState(&req.State.Session)
	action.name = req.Request.OriginalUtterance
	return h.doDeployStatus(ctx, action)
}

func (h *Handler) doDeployStatus(ctx context.Context, action *deployStatusAction) (*aliceapi.Response, errors.Err) {
	if action.namespace == "" {
		return &aliceapi.Response{
			Response: respWithTTS("В каком неймсп+ейсе проверить депл+ой?"),
			State:    action.toState(aliceapi.StateDeployStatusReqNamespace),
		}, nil
	}
	if action.namespaceID == "" {
		namespaceID, err := h.findNamespaceByName(ctx, action.namespace)
		if err != nil {
			return nil, err
		}
		if namespaceID == "" {
			respondTextF("я не нашла неймсп+ейс %s", action.namespace)
		}
		action.namespaceID = namespaceID
	}
	if action.name == "" {
		return &aliceapi.Response{
			Response: respWithTTS("Как называется депл+ой?"),
			State:    action.toState(aliceapi.StateDeployStatusReqName),
		}, nil
	}
	deploy, err := h.findDeploymentByName(ctx, action.namespaceID, action.name)
	if err != nil {
		return nil, err
	}
	if deploy == nil {
		return respondTextF("я не нашла депл+ой %s в неймсп+ейсе %s", action.name, action.namespaceID), nil
	}

	//TODO(plurals)
	if deploy.Status.UnavailableReplicas > 0 {
		return respondTextF("в депл+ое %s доступно %d реплик, еще %d в статусе анав+ейлабл",
			deploy.Name,
			deploy.Status.AvailableReplicas,
			deploy.Status.UnavailableReplicas,
		), nil
	} else {
		return respondTextF("все %d реплик в деплое %s запущены",
			deploy.Status.AvailableReplicas,
			deploy.Name,
		), nil
	}
}
