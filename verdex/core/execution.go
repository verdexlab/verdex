package core

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type Execution struct {
	Config          Config
	Targets         []string
	TargetsListPath string
	Product         string
	TestVersion     string
	OutputJsonPath  string
	Detections      []*Detection
}

func (execution *Execution) NewDetection(target string) *Detection {
	// Remove ending slash
	target = strings.TrimSuffix(target, "/")

	detection := Detection{
		Target:    target,
		Product:   execution.Product,
		Variables: make(map[string]string, 0),
	}

	// Prepend target protocol
	if !strings.HasPrefix(detection.Target, "https://") && !strings.HasPrefix(detection.Target, "http://") {
		log.Warn().Str("target", detection.Target).Msg("No protocol provided, using https by default")
		detection.Target = "https://" + detection.Target
	}

	execution.Detections = append(execution.Detections, &detection)

	return &detection
}
