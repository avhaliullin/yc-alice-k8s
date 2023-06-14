// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package stateful

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindIDByName(t *testing.T) {
	ids := []string{
		"default",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"yandex-system",
	}
	requireMatch(t, ids, "Кьюб Паблик", "kube-public")
}

func requireMatch(t *testing.T, ids []string, name string, expect string) {
	res, ok := findIdByName(name, ids)
	if !ok {
		require.Fail(t, "failed to match %s to %s", name, expect)
	}
	require.Equal(t, expect, res)
}
