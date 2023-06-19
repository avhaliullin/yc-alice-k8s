package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

type deployStatusAction struct {
	deployName  string
	deployID    string
	namespace   string
	namespaceID string
}

func deployStatusActionFromState(state *aliceapi.StateData) *deployStatusAction {
	return &deployStatusAction{
		deployName:  state.DeployName,
		deployID:    state.DeployID,
		namespace:   state.Namespace,
		namespaceID: state.NamespaceID,
	}
}

func (dsa *deployStatusAction) toState(state aliceapi.State) *aliceapi.StateData {
	return &aliceapi.StateData{
		State:       state,
		DeployName:  dsa.deployName,
		DeployID:    dsa.deployID,
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
		action.deployName = name
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
	action.deployName = req.Request.OriginalUtterance
	return h.doDeployStatus(ctx, action)
}

func (h *Handler) doDeployStatus(ctx context.Context, action *deployStatusAction) (*aliceapi.Response, errors.Err) {
	if action.namespace == "" {
		if action.deployName != "" {
			// Try in default namespace
			deploy, err := h.findDeploymentByName(ctx, k8s.DefaultNS, action.deployName)
			if err != nil {
				return nil, err
			}
			if deploy != nil {
				return resp.DeployReplicaStatuses(
					deploy.Name,
					int(deploy.Status.AvailableReplicas),
					int(deploy.Status.UnavailableReplicas),
				), nil
			}
		}
		return resp.AskNSForDeployStatus().
			WithState(action.toState(aliceapi.StateDeployStatusReqNamespace)), nil
	}
	if action.namespaceID == "" {
		namespaceID, err := h.findNamespaceByName(ctx, action.namespace)
		if err != nil {
			return nil, err
		}
		if namespaceID == "" {
			return resp.NSNotFound(action.namespace), nil
		}
		action.namespaceID = namespaceID
	}
	if action.deployName == "" {
		deployments, err := h.k8sService.ListDeployments(ctx, &k8s.ListDeploymentsReq{Namespace: action.namespaceID})
		if err != nil {
			return nil, err
		}
		if len(deployments) == 0 {
			return resp.NoDeploymentsInNS(action.namespaceID), nil
		} else if len(deployments) == 1 {
			action.deployName = deployments[0].Name
			action.deployID = action.deployName
		} else {
			var deplNames []string
			for _, deploy := range deployments {
				deplNames = append(deplNames, deploy.Name)
			}
			return resp.DeployNameForStatus(deplNames).
				WithState(action.toState(aliceapi.StateDeployStatusReqName)), nil
		}
	}
	deploy, err := h.findDeploymentByName(ctx, action.namespaceID, action.deployName)
	if err != nil {
		return nil, err
	}
	if deploy == nil {
		return resp.DeployNotFoundInNS(action.namespaceID, action.deployName), nil
	}

	return resp.DeployReplicaStatuses(
		deploy.Name,
		int(deploy.Status.AvailableReplicas),
		int(deploy.Status.UnavailableReplicas),
	), nil
}
