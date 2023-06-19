// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package resp

import (
	"fmt"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/util"
)

func respWithTTS(msg string) *aliceapi.Resp {
	tts := msg
	txt := strings.ReplaceAll(msg, "+", "")
	return &aliceapi.Resp{Text: txt, TTS: tts}
}
func respondText(msg string) *aliceapi.Response {
	return &aliceapi.Response{Response: respWithTTS(msg)}
}

func format(msg string, args ...any) respF {
	return func() *aliceapi.Response {
		return respondText(fmt.Sprintf(msg, args...))
	}
}

type respF = func() *aliceapi.Response

func randomize(opts ...respF) respF {
	if len(opts) == 0 {
		panic("at least one response option must be specified")
	}
	idx := util.RandomInt(len(opts))
	return opts[idx]
}

func concat(args ...respF) respF {
	if len(args) == 0 {
		panic("at least one part required")
	}
	if len(args) == 1 {
		return args[0]
	}
	return func() *aliceapi.Response {
		msg := ""
		for _, part := range args {
			msg += part().Response.TTS
		}
		return respondText(msg)
	}
}
