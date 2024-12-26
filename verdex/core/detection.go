package core

import (
	"time"

	"github.com/Masterminds/semver/v3"
)

type Detection struct {
	StartedAt        time.Time
	EndedAt          time.Time
	Target           string
	Product          string
	Variables        map[string]string
	TotalRequests    uint
	Success          bool
	PossibleVersions []*semver.Version
}
