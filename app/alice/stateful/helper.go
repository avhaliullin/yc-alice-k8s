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
	"k8s.io/apimachinery/pkg/util/json"
)

const maxButtons = 5

func respWithTTS(msg string) *aliceapi.Resp {
	tts := msg
	txt := strings.ReplaceAll(msg, "+", "")
	return &aliceapi.Resp{Text: txt, TTS: tts}
}
func respondText(msg string) *aliceapi.Response {
	return &aliceapi.Response{Response: respWithTTS(msg)}
}

func respondTextF(msg string, args ...any) *aliceapi.Response {
	return respondText(fmt.Sprintf(msg, args...))
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

func (h *Handler) findImageByName(ctx context.Context, imageName string) (string, errors.Err) {
	images, err := h.dockerService.ListImages(ctx)
	if err != nil {
		return "", err
	}

	result, ok := findIdByName(imageName, images)
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
	return strings.ReplaceAll(translit.ICAO(strings.ToLower(text)), " ", "-")
}

func mustToJSON(x interface{}) string {
	res, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(res)
}
