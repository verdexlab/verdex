package ui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

var (
	maxDisplayedCVEs     = 10
	maxDescriptionLength = 80
)

func RenderDetectionCVEs(execution *core.Execution, cves []*core.CVE) {
	cvesCount := len(cves)
	if cvesCount == 0 {
		log.Info().Msgf(color.New(color.BgGreen, color.FgWhite).Sprint(" ✓ ")+" Target %s (no CVE found on detected version)", color.New(color.Bold, color.FgGreen).Sprint("does not seem to be vulnerable"))
		return
	}

	vulnerableFgColor := color.FgYellow
	vulnerableBgColor := color.BgYellow
	if cves[0].CvssScore >= 7 {
		vulnerableFgColor = color.FgRed
		vulnerableBgColor = color.BgRed
	}

	log.Info().Msgf(color.New(vulnerableBgColor, color.FgWhite).Sprint(" ✗ ")+" Target is %s, %s on detected version:", color.New(color.Bold, vulnerableFgColor).Sprint("vulnerable"), color.New(color.Bold).Sprintf("%d CVE found", cvesCount))

	fmt.Println("")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("ID", "CVSS", "EPSS", "Published", "Description")
	tbl.WithHeaderFormatter(headerFmt).WithWidthFunc(tableWidthFunc)

	displayedCVEs := cves
	if cvesCount > maxDisplayedCVEs {
		displayedCVEs = cves[0:maxDisplayedCVEs]
	}

	for i, cve := range displayedCVEs {
		haveDetails := false
		description := cve.Description
		if len(description) > maxDescriptionLength {
			description = fmt.Sprintf("%s...", cve.Description[0:maxDescriptionLength])
		}

		if cve.IsKEV {
			description += fmt.Sprintf("\n • %s", color.New(color.Bold).Sprint("Known Exploited Vulnerability"))
			haveDetails = true
		}

		if haveDetails && i < cvesCount-1 {
			description += "\n"
		}

		tbl.AddRow(
			cve.ID,
			color.New(getCvssColor(cve.CvssScore)).Sprintf(" %.1f ", cve.CvssScore),
			color.New(getEpssColor(cve.EpssScore)).Sprintf("%.2f%%", cve.EpssScore*100),
			cve.PublishedAt,
			description,
		)
	}

	tbl.Print()

	if cvesCount > maxDisplayedCVEs {
		log.Info().Msgf("+ %d other CVEs", cvesCount-maxDisplayedCVEs)
	}

	fmt.Println("")

	log.Info().Msg("See output file for more information about vulnerabilities")
}

func getCvssColor(cvssScore float32) color.Attribute {
	if cvssScore >= 9 {
		return color.BgRed
	}

	if cvssScore >= 7 {
		return color.FgRed
	}

	if cvssScore >= 4 {
		return color.FgYellow
	}

	return color.FgGreen
}

func getEpssColor(epssScore float32) color.Attribute {
	if epssScore >= 0.5 {
		return color.BgRed
	}

	if epssScore >= 0.2 {
		return color.BgYellow
	}

	if epssScore >= 0.1 {
		return color.FgYellow
	}

	return color.FgCyan
}
