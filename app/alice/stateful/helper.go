package stateful

import (
	"context"
	"fmt"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/essentialkaos/translit/v2"
)

const maxButtons = 5

func respondText(msg string) *aliceapi.Response {
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: msg}}
}

func respondTextF(msg string, args ...any) *aliceapi.Response {
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: fmt.Sprintf(msg, args...)}}
}

func (h *Handler) findNamespaceByName(ctx context.Context, nsName string) (string, errors.Err) {
	nss, err := h.k8sService.ListNamespaces(ctx, &k8s.ListNamespacesReq{})
	if err != nil {
		return "", err
	}

	result, ok := findIdByName(nsName, nss)
	if !ok {
		return "", nil
	}
	return result, nil
}

func findIdByName(name string, ids []string) (string, bool) {
	idx, ok := text.BestMatch(ru2id(name), text.IDListMatcher(ids))
	if !ok {
		return "", false
	}
	return ids[idx], true
}

func ru2id(text string) string {
	return translit.ICAO(strings.ToLower(text))
}
