package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

var vertexColorBgBlackLight = color.BgRGB(54, 54, 54)
var vertexColorBgBlackDark = color.BgRGB(42, 42, 42)

var vertexColorBgGreenLight = color.BgRGB(78, 180, 146)
var vertexColorBgGreenDark = color.BgRGB(62, 141, 112)

var vertexColorBgRedLight = color.BgRGB(180, 78, 78)
var vertexColorBgRedDark = color.BgRGB(141, 62, 62)

var vertexColorFgGreenLight = color.RGB(78, 180, 146)
var vertexColorFgRedLight = color.RGB(180, 78, 78)

type UiDetectionVertex struct {
	versions  []*semver.Version
	colorsMap []bool
	header    string
}

func NewDetectionVertex(versions []*semver.Version) *UiDetectionVertex {
	vertex := UiDetectionVertex{
		versions:  versions,
		colorsMap: make([]bool, len(versions)),
		header:    "",
	}

	isColorLight := false

	for index, version := range vertex.versions {
		isNewMajorVersion := index == 0 || vertex.versions[index-1].Major() != version.Major()

		if isNewMajorVersion {
			// Toggle color intensity (light or dark) at each major version
			isColorLight = !isColorLight

			majorString := fmt.Sprintf("%d", version.Major())

			// Determine header: if we should display major version number, or spaces because not enough space
			displayedString := majorString
			for nextVersionIndex, nextVersion := range vertex.versions[index+1 : len(vertex.versions)] {
				// loop until next major version
				if nextVersion.Major() != version.Major() {
					// if not enough space to display major version number, display nothing
					if len(majorString) > nextVersionIndex {
						displayedString = strings.Repeat(" ", nextVersionIndex+1)
					}

					break
				}

				// add spaces suffix
				if nextVersionIndex >= len(majorString)-1 {
					displayedString += " "
				}
			}

			vertex.header += displayedString
		}

		vertex.colorsMap[index] = isColorLight
	}

	return &vertex
}

// Render vertex header with major versions numbers
func (vertex *UiDetectionVertex) RenderHeader() {
	fmt.Println("")
	log.Info().Msg(color.HiBlackString(vertex.header))
}

// Render a detected rule line
func (vertex *UiDetectionVertex) RenderMatchingLine(constraint *semver.Constraints) {
	row := ""

	for index, version := range vertex.versions {
		isMatch, _ := constraint.Validate(version)

		var c *color.Color
		if isMatch {
			if vertex.colorsMap[index] {
				c = vertexColorBgGreenLight
			} else {
				c = vertexColorBgGreenDark
			}
		} else {
			if vertex.colorsMap[index] {
				c = vertexColorBgBlackLight
			} else {
				c = vertexColorBgBlackDark
			}
		}

		row += c.Sprint(" ")
	}

	row += "  " + vertexColorFgGreenLight.Sprint(constraint.String())

	log.Info().Msg(row)
}

// Render all excluded rules on single line (bottom)
func (vertex *UiDetectionVertex) RenderExcludedLine(excludedVersions []*semver.Version) {
	row := ""

	excludedVersionsStr := make([]string, len(excludedVersions))
	for index, version := range excludedVersions {
		excludedVersionsStr[index] = version.String()
	}

	for index, version := range vertex.versions {
		var c *color.Color

		if slices.Contains(excludedVersionsStr, version.String()) {
			if vertex.colorsMap[index] {
				c = vertexColorBgRedLight
			} else {
				c = vertexColorBgRedDark
			}
		} else {
			if vertex.colorsMap[index] {
				c = vertexColorBgBlackLight
			} else {
				c = vertexColorBgBlackDark
			}
		}

		row += c.Sprint(" ")
	}

	row += "  " + vertexColorFgRedLight.Sprint(strings.Join(excludedVersionsStr, ", "))

	log.Info().Msg(row)
}

// Render footer, with possible versions
func (vertex *UiDetectionVertex) RenderPossibleVersions(possibleVersions []*semver.Version) {
	if len(possibleVersions) == 0 {
		log.Info().Msg(color.RedString(" → no matching rule found\n"))
		return
	}

	row := ""

	possibleVersionsStr := make([]string, len(possibleVersions))
	for index, version := range possibleVersions {
		possibleVersionsStr[index] = version.String()
	}

	possibleVersionColor := vertexColorFgGreenLight
	if len(possibleVersions) > 1 {
		possibleVersionColor = color.New(color.FgYellow)
	}

	for _, version := range vertex.versions {
		isPossibleVersion := slices.Contains(possibleVersionsStr, version.String())

		if isPossibleVersion {
			row = row + possibleVersionColor.Sprint("↑")
		} else {
			row = row + " "
		}
	}

	row += "\n"
	log.Info().Msg(row)
}
