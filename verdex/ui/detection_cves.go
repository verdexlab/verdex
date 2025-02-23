package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/api"
	"github.com/verdexlab/verdex/verdex/core"
)

func RenderDetectionCVEs(execution *core.Execution, data *api.CVEsData, isApiKeyValid bool) {
	fmt.Println("")

	if !data.Vulnerable {
		log.Info().Msgf(color.New(color.BgGreen, color.FgWhite).Sprint(" ✓ ")+" Target is %s (no registered CVE on detected version)", color.New(color.Bold, color.FgGreen).Sprint("does not seem to be vulnerable"))
		return
	}

	log.Info().Msgf(color.New(color.BgRed, color.FgWhite).Sprint(" ✗ ")+" Target is %s, %s on detected version", color.New(color.Bold, color.FgRed).Sprint("vulnerable"), color.New(color.Bold).Sprintf("%d CVE found", data.Total))
	if !isApiKeyValid {
		log.Info().Msg("Use valid API key (-key) to list vulnerabilities and get update recommendations, see https://docs.verdexlab.io/expert/vulnerabilities")
		return
	}

	fmt.Println("")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("ID", "CVSS", "EPSS", "Published", "NPV*", "Summary")
	tbl.WithHeaderFormatter(headerFmt).WithWidthFunc(tableWidthFunc)

	for i, cve := range data.CVEs {
		summary := cve.Summary
		haveDetails := false

		if cve.KevSince != nil {
			summary += fmt.Sprintf("\n • %s since %s", color.New(color.Bold).Sprint("Known Exploited Vulnerability"), *cve.KevSince)
			haveDetails = true
		}

		for _, poc := range cve.POCs {
			summary += fmt.Sprintf("\n • %s %s", color.New(color.Bold).Sprint("POC:"), poc)
			haveDetails = true
		}

		if cve.NucleiTemplate != nil {
			summary += fmt.Sprintf("\n • %s %s", color.New(color.Bold).Sprint("Nuclei Template:"), *cve.NucleiTemplate)
			haveDetails = true
		}

		if haveDetails && i < len(data.CVEs)-1 {
			summary += "\n"
		}

		if cve.NearestPatchedVersion == nil {
			dash := "-"
			cve.NearestPatchedVersion = &dash
		}

		tbl.AddRow(
			cve.ID,
			color.New(getCvssColor(cve.CvssScore)).Sprintf(" %.1f ", cve.CvssScore),
			color.New(getEpssColor(cve.EpssScore)).Sprintf("%.2f%%", cve.EpssScore*100),
			strings.Split(cve.PublishedAt, "T")[0],
			*cve.NearestPatchedVersion,
			summary,
		)
	}

	tbl.Print()

	fmt.Println("")

	if data.UpdateRecommendations != nil && (data.UpdateRecommendations.WithoutVulnerabilities != nil || data.UpdateRecommendations.WithoutCriticalVulnerabilities != nil) {
		log.Info().Msg(color.New(color.Bold).Sprint("Update recommendations:"))

		msg := " • Nearest version without vulnerabilities: "
		if data.UpdateRecommendations.WithoutVulnerabilities != nil {
			msg += color.New(color.Bold, color.FgGreen).Sprint(*data.UpdateRecommendations.WithoutVulnerabilities)
		} else {
			msg += color.New(color.Bold).Sprint("N/A")
		}
		log.Info().Msg(msg)

		msg = " • Nearest version without critical vulnerabilities (CVSS >= 9): "
		if data.UpdateRecommendations.WithoutCriticalVulnerabilities != nil {
			msg += color.New(color.Bold, color.FgYellow).Sprint(*data.UpdateRecommendations.WithoutCriticalVulnerabilities)
		} else {
			msg += color.New(color.Bold).Sprint("N/A")
		}
		log.Info().Msg(msg)
	} else {
		log.Info().Msg("Cannot determine update recommendations for these versions")
	}

	fmt.Println("")
	log.Info().Msg("See output file for more information about vulnerabilities")
	log.Info().Msgf("%s Nearest future version not vulnerable to given CVE", color.New(color.Bold).Sprint("*NPV (Nearest Patched Version):"))
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
