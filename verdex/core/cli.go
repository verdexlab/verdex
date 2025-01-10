package core

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Parse CLI arguments
func ParseFlags() *Execution {
	execution := Execution{
		Config: Config{},
	}

	flag.Usage = func() {
		fmt.Println("Verdex is a fast detection tool for technologies versions.\n")
		fmt.Println("Usage:\n  verdex [flags]\n")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	var showVersion bool
	var singleTarget string

	// Targets
	flag.StringVar(&singleTarget, "target", "", "target to scan")
	flag.StringVar(&execution.TargetsListPath, "list", "", "path to list of targets to scan")
	flag.StringVar(&execution.Product, "product", "", "specify product to scan")

	// Outputs
	flag.StringVar(&execution.OutputJsonPath, "output-json", "", "file path to output results in JSON format")

	// Misc
	flag.StringVar(&execution.Config.ApiKey, "key", "", "provide Verdex API key")
	flag.BoolVar(&showVersion, "version", false, "show verdex version")
	flag.StringVar(&execution.Config.TemplatesDirectory, "templates-directory", "", "path to templates directory to use instead of official repository")
	flag.BoolVar(&execution.Config.Verbose, "verbose", false, "show verbose output")
	flag.BoolVar(&execution.Config.ReportTargets, "report-errors", false, "report failed target to improve detections")

	// Testing
	if GetEnvironment() == EnvironmentDevelopment {
		flag.BoolVar(&execution.Config.Test, "test", false, "run unit tests")
		flag.StringVar(&execution.Config.TestVersion, "test-version", "", "test only a specific version")
		flag.BoolVar(&execution.Config.TestSession, "test-session", false, "start a test session for given version")
	}

	flag.Parse()

	if showVersion {
		log.Info().Msgf("Verdex Version: %s", GetVerdexVersion())
		os.Exit(0)
	}

	if singleTarget == "" && execution.TargetsListPath == "" && !execution.Config.Test {
		log.Fatal().Msg("either -target or -list argument is required (use `verdex -help` for more information)")
	} else if singleTarget != "" && execution.TargetsListPath != "" && !execution.Config.Test {
		log.Fatal().Msg("only one of the -target and -list arguments is allowed (use `verdex -help` for more information)")
	}

	// Parse single target
	if singleTarget != "" {
		execution.Targets = append(execution.Targets, singleTarget)
	}

	// Parse multiple targets
	if execution.TargetsListPath != "" {
		targetsList, err := os.ReadFile(execution.TargetsListPath)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to read -list file: %s", execution.TargetsListPath)
		}

		targets := strings.Split(string(targetsList), "\n")

		for _, target := range targets {
			target = strings.Trim(target, " ")
			if target == "" {
				continue
			}

			execution.Targets = append(execution.Targets, target)
		}

		if len(execution.Targets) == 0 {
			log.Fatal().Msg("-list file is empty")
		}
	}

	// Parse product
	if execution.Config.Test && execution.Product == "" {
		log.Fatal().Msg("-product argument is required with -test (use `verdex -help` for more information)")
	}

	if execution.Config.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Set local templates directory for local development run
	if GetEnvironment() == EnvironmentDevelopment && execution.Config.TemplatesDirectory == "" {
		log.Warn().Msg("Using local templates directory for development")
		execution.Config.TemplatesDirectory = "./templates"
	}

	if execution.Config.TemplatesDirectory == "" {
		execution.Config.TemplatesSource = TemplatesSourceGitHubOfficial
		execution.Config.TemplatesOrganization = TemplatesOfficialOrganization
		execution.Config.TemplatesRepository = TemplatesOfficialRepository
		execution.Config.TemplatesDirectory = TemplatesDefaultDirectory
	} else {
		execution.Config.TemplatesSource = TemplatesSourceLocalDirectory
	}

	// Test session
	if execution.Config.TestSession && execution.Config.TestVersion == "" {
		log.Fatal().Msg("-test-version argument is required with -test-session (use `verdex -help` for more information)")
	}

	return &execution
}
