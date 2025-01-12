package core

import (
	"time"

	"github.com/Masterminds/semver/v3"
)

type Detection struct {
	StartedAt             time.Time
	EndedAt               time.Time
	Target                string
	Product               string
	Variables             map[string]string
	TotalRequests         uint
	Success               bool
	PossibleVersions      []*semver.Version
	Vulnerable            bool
	CVEs                  []CVE
	UpdateRecommendations *DetectionUpdateRecommendations `json:"update_recommendations,omitempty"`
}

type DetectionUpdateRecommendations struct {
	WithoutVulnerabilities         *string `json:"without_vulnerabilities"`
	WithoutCriticalVulnerabilities *string `json:"without_critical_vulnerabilities"`
}
