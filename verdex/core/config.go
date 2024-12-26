package core

import (
	"os"
	"path"
)

type Config struct {
	TemplatesSource       TemplatesSource
	TemplatesOrganization string
	TemplatesRepository   string
	TemplatesDirectory    string
	TemplatesRelease      string
	Verbose               bool
	ReportTargets         bool
	Test                  bool
	TestVersion           string
	TestSession           bool
}

type TemplatesSource string

const (
	TemplatesSourceGitHubOfficial TemplatesSource = "github-official"
	TemplatesSourceLocalDirectory TemplatesSource = "local-directory"
)

var userHomeDir, _ = os.UserHomeDir()

// CLI version
var cliVersion = "0.0.1"

// Templates
var TemplatesDefaultDirectory = path.Join(userHomeDir, "verdex-templates")
var TemplatesOfficialOrganization = "verdexlab"
var TemplatesOfficialRepository = "verdex"
