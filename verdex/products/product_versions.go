package products

import (
	"github.com/Masterminds/semver/v3"
)

type ProductVersions struct {
	Source       ProductVersionsSource `yaml:"source" validate:"required,oneof=list github"`
	List         []*semver.Version     `yaml:"-"`
	RawList      []string              `yaml:"list" validate:"required_if=Source list"`
	Organization string                `yaml:"organization" validate:"required_if=Source github"`
	Repository   string                `yaml:"repository" validate:"required_if=Source github"`
	cachePath    string                `yaml:"-"`
}

type ProductVersionsSource string

const (
	ProductVersionsSourceList   ProductVersionsSource = "list"
	ProductVersionsSourceGitHub ProductVersionsSource = "github"
)

// Get all product versions that are matching included & excluded constraints
func (productVersions *ProductVersions) GetVersionsMatchingConstraints(included []*semver.Constraints, excluded []*semver.Constraints) (m []*semver.Version, e []*semver.Version, err error) {
	matchingVersions := make([]*semver.Version, 0)
	excludedVersions := make([]*semver.Version, 0)

	for _, version := range productVersions.List {
		isMatching := true
		isExcluded := false

		for _, c := range included {
			isIncluded, _ := c.Validate(version)
			if !isIncluded {
				isMatching = false
				break
			}
		}

		if isMatching {
			for _, c := range excluded {
				isNotIncluded, _ := c.Validate(version)
				if isNotIncluded {
					isMatching = false
					isExcluded = true
					break
				}
			}
		}

		if isMatching {
			matchingVersions = append(matchingVersions, version)
		} else if isExcluded {
			excludedVersions = append(excludedVersions, version)
		}
	}

	return matchingVersions, excludedVersions, nil
}
