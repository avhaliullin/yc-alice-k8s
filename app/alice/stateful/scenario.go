package stateful

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
)

type scenario = func(context.Context, *aliceapi.Request) (*aliceapi.Response, errors.Err)

func (h *Handler) setupScenarios() {
	h.stateScenarios = map[aliceapi.State]scenario{
		aliceapi.StateDeployReqName:            h.deployReqName,
		aliceapi.StateDeployReqImage:           h.deployReqImage,
		aliceapi.StateDeployConfirm:            h.deployReqConfirm,
		aliceapi.StateDeployStatusReqName:      h.deployStatusReqName,
		aliceapi.StateDeployStatusReqNamespace: h.deployStatusReqNamespace,
		aliceapi.StateScaleDeployReqName:       h.scaleDeployReqName,
		aliceapi.StateScaleDeployReqScale:      h.scaleDeployReqScale,
		aliceapi.StateScaleDeployReqConfirm:    h.scaleDeployReqConfirm,
		aliceapi.StateDeleteDeployReqName:      h.deleteDeployReqName,
		aliceapi.StateDeleteDeployReqConfirm:   h.deleteDeployReqConfirm,
		aliceapi.StateBrokenPodsReqNS:          h.brokenPodsReqNs,
		aliceapi.StateCountPodsReqNS:           h.countPodsReqNs,
		aliceapi.StateIngressListReqNS:         h.listIngressesReqNs,
		aliceapi.StateServiceListReqNS:         h.listServicesReqNs,
	}
	h.scratchScenarios = []scenario{
		h.listNamespaces,
		h.countPods,
		h.brokenPods,
		h.listServices,
		h.listIngresses,
		h.discoverScenarios,
		h.easterEggs,
		h.deployFromScratch,
		h.deployStatusFromScratch,
		h.scaleDeployFromScratch,
		h.deleteDeployFromScratch,
		h.letsPlayK8S,
	}
}
