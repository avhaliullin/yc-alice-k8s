package alice

import (
	"context"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
)

type Handler interface {
	Handle(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, error)
}
