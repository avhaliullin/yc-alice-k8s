package stateful

import (
	"context"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/cache"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/text/resp"
	docker_hub "github.com/avhaliullin/yandex-alice-k8s-skill/app/docker-hub"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/errors"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/log"
	"go.uber.org/zap"
)

type Handler struct {
	logger           *zap.Logger
	stateScenarios   map[aliceapi.State]scenario
	scratchScenarios []scenario
	k8sService       k8s.Service
	dockerService    docker_hub.Service
}

func NewHandler(deps Deps) (*Handler, error) {
	h := &Handler{
		logger:        deps.GetLogger(),
		k8sService:    deps.GetK8sService(),
		dockerService: deps.GetDockerService(),
	}
	h.setupScenarios()
	return h, nil
}

type intentWa struct {
	phrase string
	intent aliceapi.Intents
}

var intents = intentsMatcher([]intentWa{
	{phrase: "сколько подов", intent: aliceapi.Intents{CountPods: &aliceapi.IntentCountPods{}}},
	{phrase: "посчитай поды", intent: aliceapi.Intents{CountPods: &aliceapi.IntentCountPods{}}},
	{phrase: "отмасштабируй деплой", intent: aliceapi.Intents{ScaleDeploy: &aliceapi.IntentScaleDeploy{}}},
	{phrase: "как ты работаешь", intent: aliceapi.Intents{HowYouMade: &aliceapi.EmptyObj{}}},
	{phrase: "как ты сделана", intent: aliceapi.Intents{HowYouMade: &aliceapi.EmptyObj{}}},
	{phrase: "можно ли запустить базу данных в кубере", intent: aliceapi.Intents{EasterDBLaunch: &aliceapi.EmptyObj{}}},
	{phrase: "можно ли задеплоить базу данных в кубер", intent: aliceapi.Intents{EasterDBLaunch: &aliceapi.EmptyObj{}}},
})

var _ text.MatchCandidates = intentsMatcher([]intentWa{})

type intentsMatcher []intentWa

func (i intentsMatcher) Len() int {
	return len(i)
}

func (i intentsMatcher) TextOf(idx int) string {
	return i[idx].phrase
}

func (h *Handler) matchNoIntent(req *aliceapi.Request) *aliceapi.Intents {
	res, ok := text.BestMatch(strings.ToLower(req.Request.OriginalUtterance), intents, text.MatchMinRatio(0.9))
	if !ok {
		return nil
	}
	return &intents[res].intent
}

func (h *Handler) Handle(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, error) {
	sessionID := req.Session.SessionID
	ctx = log.CtxWithLogger(ctx, h.logger.With(zap.String("sessionID", string(sessionID))))
	ctx = cache.ContextWithCache(ctx)
	resp, err := h.handle(ctx, req)
	if err != nil {
		return h.reportError(ctx, err)
	}
	resp.Version = req.Version
	log.Info(ctx, "request processed",
		log.FieldJSON("req", req),
		log.FieldJSON("resp", resp),
		zap.String("kind", "access-log"),
	)
	return resp, nil
}

func (h *Handler) handle(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Session.New || req.AccountLinkingComplete != nil {
		return resp.WelcomePhrase(), nil
	}
	if state := req.State.Session; state.State != aliceapi.StateInit {
		intents := req.Request.NLU.Intents
		if req.Request.Type == aliceapi.RequestTypeSimple && intents.Cancel != nil {
			return resp.RejectOnWizard(), nil
		}
		scenario, ok := h.stateScenarios[state.State]
		if ok {
			resp, err := scenario(ctx, req)
			if err != nil {
				return nil, err
			}
			if resp != nil {
				return resp, nil
			}
		}
	}
	for _, s := range h.scratchScenarios {
		resp, err := s(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp != nil {
			return resp, err
		}
	}
	parsedIntent := h.matchNoIntent(req)
	if parsedIntent != nil {
		req.Request.NLU.Intents = *parsedIntent
		for _, s := range h.scratchScenarios {
			resp, err := s(ctx, req)
			if err != nil {
				return nil, err
			}
			if resp != nil {
				return resp, err
			}
		}
	}
	return resp.UnrecognizedRequest(), nil
}

func (h *Handler) reportError(ctx context.Context, err errors.Err) (*aliceapi.Response, error) {
	errors.Log(ctx, err)
	return nil, err
}
