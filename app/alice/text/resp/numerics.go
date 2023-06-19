// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package resp

import (
	"fmt"
	"strconv"
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

func nvl(a, b respF) respF {
	if a != nil {
		return a
	} else {
		return b
	}
}

func (p *numDependentConfig) get(num int) respF {
	switch num {
	case 0:
		return nvl(p.exactly0, p.like5)
	case 1:
		return nvl(p.exactly1, p.like1)
	case 2:
		return nvl(p.exactly2, p.like2)
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

func number(val int, gCase GramCase, gGender GramGender) string {
	strVal := strconv.Itoa(val)
	if gCase == CaseNominative && gGender == GenderM {
		return strVal
	}
	twoDigits := val % 100
	oneDigit := val % 10
	if gCase == CaseAccusative && gGender == GenderF {
		if twoDigits >= 10 && twoDigits < 20 || (oneDigit != 1 && oneDigit != 2) {
			return strVal
		}
		pref := ""
		if oneDigit != val {
			pref = strconv.Itoa(val-oneDigit) + " "
		}
		switch oneDigit {
		case 1:
			return pref + "одну"
		case 2:
			return pref + "две"
		}
	}
	return strVal
}

type GramCase int

const (
	// один под
	CaseNominative GramCase = iota
	// одного пода
	CaseAccusative
)

type GramGender int

const (
	GenderF GramGender = iota
	GenderM
	GenderN
)
