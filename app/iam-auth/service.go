// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package iam_auth

import (
	"context"
)

type Service interface {
	GetToken(ctx context.Context) (string, error)
}
