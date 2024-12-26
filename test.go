package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/detect"
	"github.com/verdexlab/verdex/verdex/products"
	"github.com/verdexlab/verdex/verdex/templates"
	"github.com/verdexlab/verdex/verdex/tests"
)

var maxSimultaneousTestContainers = 3

func RunTests(execution *core.Execution, product *products.Product) {
	err := product.Versions.ReloadList()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load product versions list")
	}

	log.Info().Str("product", product.ID).Msg("Starting Verdex tests")

	testCases := tests.GetProductTestCases(product.ID)
	if len(testCases) == 0 {
		log.Warn().Msgf("[%s] No test cases found!", product.ID)
		return
	}

	spinnerManager := ysmrr.NewSpinnerManager(
		ysmrr.WithSpinnerColor(colors.FgHiBlue),
	)
	spinnerManager.Start()

	var wg sync.WaitGroup

	for _, testCase := range testCases {
		versionsConstraints := []*semver.Constraints{testCase.Constraint}
		if execution.Config.TestVersion != "" {
			constraint, err := semver.NewConstraint(execution.Config.TestVersion)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid -test-version argument (for syntax, see https://github.com/Masterminds/semver)")
			}

			versionsConstraints = append(versionsConstraints, constraint)
		}

		versionsToTest, _, err := product.Versions.GetVersionsMatchingConstraints(versionsConstraints, []*semver.Constraints{})
		if err != nil {
			spinnerManager.Stop()
			log.Fatal().Err(err).Msg("Failed to get versions matching constraints")
		}

		if len(versionsToTest) == 0 {
			log.Warn().Msgf("[%s/%s] No versions to test found", testCase.Info.Product, testCase.Name)
			continue
		}

		rounds := len(versionsToTest) / maxSimultaneousTestContainers
		if rounds*maxSimultaneousTestContainers < len(versionsToTest) {
			rounds++
		}

		for round := 0; round < rounds; round++ {
			for _, version := range versionsToTest[round*maxSimultaneousTestContainers : min(len(versionsToTest), (round+1)*maxSimultaneousTestContainers)] {
				wg.Add(1)
				go runTestCaseOnVersion(execution, testCase, version, spinnerManager, &wg)
				time.Sleep(100 * time.Millisecond) // preserve order
			}

			wg.Wait()
		}
	}

	spinnerManager.Stop()
}

func runTestCaseOnVersion(execution *core.Execution, testCase *tests.TestCase, version *semver.Version, spinnerManager ysmrr.SpinnerManager, wg *sync.WaitGroup) {
	defer wg.Done()

	spinner := spinnerManager.AddSpinner("initializing instance")

	setPrefixColor := func(c color.Attribute) {
		spinner.UpdatePrefixf("%s %s ", testCase.Info.Product, color.New(c).Sprint(version))
	}
	setPrefixColor(color.FgCyan)

	instance, err := testCase.NewInstance(version.String())
	if err != nil {
		setPrefixColor(color.FgRed)
		spinner.ErrorWithMessage("failed")
		log.Error().Err(err).Msgf("[%s/%s/%s] Failed to initialize test instance", testCase.Info.Product, testCase.Name, version)
		instance.Destroy()
		return
	}

	spinner.UpdateMessage("pulling docker image")
	err = instance.Up()
	if err != nil {
		setPrefixColor(color.FgRed)
		spinner.ErrorWithMessage("failed")
		log.Error().Err(err).Msgf("[%s/%s/%s] Failed to up test instance", testCase.Info.Product, testCase.Name, version)
		instance.Destroy()
		return
	}

	spinner.UpdateMessage("starting instance")
	for {
		time.Sleep(5 * time.Second)

		healthy, err := instance.ServiceIsHealthy()
		if err != nil {
			setPrefixColor(color.FgRed)
			spinner.ErrorWithMessage("failed")
			log.Error().Err(err).Msgf("[%s/%s/%s] Failed to up test instance", testCase.Info.Product, testCase.Name, version)
			instance.Destroy()
			return
		} else if healthy {
			break
		}
	}

	target, err := instance.ServiceOrigin()
	if err != nil {
		spinner.ErrorWithMessage("failed")
		log.Error().Err(err).Msgf("[%s/%s/%s] Failed to determine service origin", testCase.Info.Product, testCase.Name, version)
		instance.Destroy()
		return
	}

	if execution.Config.TestSession {
		for {
			spinner.UpdateMessagef("at %s | [d]etect or [s]top", target)

			var action []byte = make([]byte, 1)
			os.Stdin.Read(action)

			if string(action) == "d" {
				if execution.Config.TemplatesSource == core.TemplatesSourceLocalDirectory {
					spinner.UpdateMessage("reloading templates")
					templates.LoadTemplatesFromDirRecursively(&execution.Config)
				}

				spinner.UpdateMessage("running version detection")
				completeMessage, errorMessage, _ := versionDetectionOnVersion(execution, version, target)
				if completeMessage != "" {
					log.Info().Msgf("%s %s : %s", execution.Product, version, completeMessage)
				} else {
					log.Error().Msgf("%s %s : %s", execution.Product, version, errorMessage)
				}
			} else if string(action) == "s" {
				spinner.UpdateMessage("destroying instance")
				instance.Destroy()
				spinner.CompleteWithMessage("session done")
				return
			}
		}
	}

	spinner.UpdateMessage("running version detection")

	completeMessage, errorMessage, messageColor := versionDetectionOnVersion(execution, version, target)
	setPrefixColor(messageColor)
	if completeMessage != "" {
		spinner.CompleteWithMessage(completeMessage)
	} else {
		spinner.ErrorWithMessage(errorMessage)
	}

	instance.Destroy()
}

func versionDetectionOnVersion(execution *core.Execution, version *semver.Version, target string) (completeMessage string, errorMessage string, messageColor color.Attribute) {
	detection := core.Detection{
		Target:    target,
		Product:   execution.Product,
		Variables: make(map[string]string, 0),
	}

	versions, err := detect.DetectVersion(execution, &detection)
	if err != nil {
		messageColor = color.FgRed
		errorMessage = fmt.Sprintf("error: %s", err)
		return
	}

	if len(versions) == 1 && versions[0].Equal(version) {
		messageColor = color.FgGreen
		completeMessage = "OK"
		return
	}

	if len(versions) == 0 {
		messageColor = color.FgRed
		errorMessage = "failed: no version"
		return
	}

	versionsResults := make([]string, 0)
	for _, detectedVersion := range versions {
		if detectedVersion.Equal(version) {
			versionsResults = append(versionsResults, color.New(color.FgGreen).Sprintf("%s ✓", detectedVersion))
		} else {
			versionsResults = append(versionsResults, color.New(color.FgRed).Sprintf("%s ✗", detectedVersion))
		}
	}

	var versionIsInaccurate bool

	for _, v := range versions {
		if v.Equal(version) {
			versionIsInaccurate = true
			break
		}
	}

	message := strings.Join(versionsResults, color.HiBlackString(" | "))

	if versionIsInaccurate {
		messageColor = color.FgYellow
		completeMessage = "inaccurate: " + message
	} else {
		messageColor = color.FgRed
		errorMessage = "failed: " + message
	}

	return
}
