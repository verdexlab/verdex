package rules

import "github.com/verdexlab/verdex/verdex/assets"

type RuleHttpMatcher struct {
	Type   RuleHttpMatcherType `yaml:"type" validate:"required,oneof=word regex status"`
	Part   RuleHttpMatcherPart `yaml:"part" validate:"required_if=Type word,required_if=Type regex,len=0|oneof=body"`
	Word   string              `yaml:"word" validate:"required_if=Type word"`
	Regex  string              `yaml:"regex" validate:"required_if=Type regex"`
	Status int                 `yaml:"status" validate:"required_if=Type status,gte=0,lte=999"`
}

type RuleHttpMatcherType string

const (
	RuleHttpMatcherTypeWord   RuleHttpMatcherType = "word"
	RuleHttpMatcherTypeRegex  RuleHttpMatcherType = "regex"
	RuleHttpMatcherTypeStatus RuleHttpMatcherType = "status"
)

type RuleHttpMatcherPart string

const (
	RuleHttpMatcherPartBody RuleHttpMatcherPart = "body"
)

func (matcher *RuleHttpMatcher) Match(asset *assets.Asset) bool {
	if matcher.Type == RuleHttpMatcherTypeWord {
		return asset.BodyMatchWord(matcher.Word)
	}

	if matcher.Type == RuleHttpMatcherTypeRegex {
		return asset.BodyMatchRegex(matcher.Regex)
	}

	if matcher.Type == RuleHttpMatcherTypeStatus {
		return matcher.Status == asset.StatusCode
	}

	return false
}
