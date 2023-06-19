// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package resp

import (
	"fmt"
)

func numDependent(num int, conf numDependentConfig) respF {
	return conf.get(num)
}

type numDependentConfig struct {
	exactly0 respF
	exactly1 respF
	exactly2 respF
	like1    respF
	like2    respF
	like5    respF
}

func (p *numDependentConfig) get(num int) respF {
	switch num {
	case 0:
		return p.exactly0
	case 1:
		return p.exactly1
	case 2:
		return p.exactly2
	}
	if num < 0 {
		num = -num
	}
	num = num % 100
	if num >= 10 && num < 20 {
		return p.like5
	}
	num = num % 10
	switch num {
	case 0, 5, 6, 7, 8, 9:
		return p.like5
	case 1:
		return p.like1
	case 2, 3, 4:
		return p.like2
	default:
		panic(fmt.Sprintf("uncovered num: %d", num))
	}
}

/*
0, 5, 6, 7, 8, 9, 10, ... 20, 200, 2000 арбузов
1 арбуз
2, 3, 4 арбуза

*/
