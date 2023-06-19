package text

import (
	"strings"
	"unicode"

	"github.com/texttheater/golang-levenshtein/levenshtein"
	appsv1 "k8s.io/api/apps/v1"
)

var levensteinOpts = levenshtein.Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 1,
	Matches: func(r rune, r2 rune) bool {
		return unicode.ToLower(r) == unicode.ToLower(r2)
	},
}

type options struct {
	minRatio float64
	prefix   string
}

type MatchCandidates interface {
	Len() int
	TextOf(idx int) string
}

func BestMatch(text string, candidates MatchCandidates, opts ...MatchOpt) (int, bool) {
	options := options{minRatio: 0.5}
	for _, opt := range opts {
		opt(&options)
	}
	n := candidates.Len()
	if n == 0 {
		return -1, false
	}
	bestRatio := 0.0
	bestMatchIdx := -1

	matchText := func(idx int, text string, candidate string) {
		ratio := levenshtein.RatioForStrings([]rune(text), []rune(candidate), levensteinOpts)
		if ratio > bestRatio {
			bestRatio = ratio
			bestMatchIdx = idx
		}
	}
	for idx := 0; idx < n; idx++ {
		candidateText := candidates.TextOf(idx)
		matchText(idx, text, candidateText)
		if options.prefix != "" {
			matchText(idx, options.prefix+" "+text, candidateText)
			matchText(idx, text, options.prefix+" "+candidateText)
		}
	}
	if bestRatio < options.minRatio {
		return -1, false
	}
	return bestMatchIdx, true
}

type MatchOpt func(*options)

func MatchMinRatio(r float64) MatchOpt {
	return func(o *options) {
		o.minRatio = r
	}
}

func MatchOptPrefix(prefix string) MatchOpt {
	return func(o *options) {
		o.prefix = prefix
	}
}

var _ MatchCandidates = IDListMatcher([]string{})

type IDListMatcher []string

func (m IDListMatcher) Len() int {
	return len(m)
}

func (m IDListMatcher) TextOf(idx int) string {
	return normalize(m[idx])
}

var _ MatchCandidates = DeploymentsMatcher([]appsv1.Deployment{})

type DeploymentsMatcher []appsv1.Deployment

func (d DeploymentsMatcher) Len() int {
	return len(d)
}

func (d DeploymentsMatcher) TextOf(idx int) string {
	return normalize(d[idx].Name)
}

func normalize(id string) string {
	return strings.ToLower(id)
}
