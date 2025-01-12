package templates

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v30/github"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/thirdparty"
	"github.com/verdexlab/verdex/verdex/ui"
)

var downloadDirectory = "templates"

var templatesReleasesPrefix = "templates-"

func UpdateLatestRelease(config *core.Config) error {
	release, data, err := downloadLatestRelease(config)
	if err != nil {
		return err
	}

	deleteAndRecreateTemplatesDirectory(config)

	err = unpackReleaseWithCallback(config.TemplatesDirectory, data)
	if err != nil {
		return err
	}

	cache := core.GetCache(config)
	cache.Releases.Templates.Current = *(release.Name)
	cache.Save()

	return nil
}

/**
 * Download latest release of GitHub repository
 */
func downloadLatestRelease(config *core.Config) (*github.RepositoryRelease, *bytes.Reader, error) {
	githubClient, httpClient := thirdparty.GitHubGetClients()

	release, err := thirdparty.GitHubGetLatestPrefixedRelease(config.TemplatesOrganization, config.TemplatesRepository, templatesReleasesPrefix, githubClient)
	if err != nil {
		return nil, nil, err
	}

	downloadURL := release.GetZipballURL()
	res, err := httpClient.Get(downloadURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to source GitHub zipball URL: %s", downloadURL)
	}

	defer res.Body.Close()

	// progress bar
	bar := ui.ProgressBarStart(res.Body, res.ContentLength)
	defer ui.ProgressBarFinish(bar)

	bin, err := io.ReadAll(bar.Reader)
	if err != nil {
		return nil, nil, errors.New("failed to read GitHub zipball body")
	}

	return release, bytes.NewReader(bin), nil
}

/**
 * Delete given directory if exists, and recreate it
 */
func deleteAndRecreateTemplatesDirectory(config *core.Config) {
	// Backup cache content
	cache := core.GetCache(config)

	// Delete previous directory if already exists
	if _, err := os.Stat(config.TemplatesDirectory); !os.IsNotExist(err) {
		os.RemoveAll(config.TemplatesDirectory)
	}

	// Recreate directory
	os.MkdirAll(config.TemplatesDirectory, os.ModePerm)

	// Restore cache content
	cache.Save()
}

/**
 * Unpack GitHub release and write yml files to destination directory
 */
func unpackReleaseWithCallback(destDir string, data *bytes.Reader) error {
	callbackFunc := func(uri string, f fs.FileInfo, r io.Reader) error {
		uriParts := strings.Split(uri, "/")

		// example: verdexlab-verdex-a0b1c2d3/templates/keycloak/rules/26.0.5.yml
		if len(uriParts) < 2 || uriParts[1] != downloadDirectory || f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(uri, ".yml") && !strings.HasSuffix(uri, ".yaml") {
			return nil
		}

		writeDirectory := path.Join(destDir, strings.Join(uriParts[2:len(uriParts)-1], "/"))
		writeFilePath := path.Join(writeDirectory, uriParts[len(uriParts)-1])

		bin, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read release file %s", uri)
		}

		err = os.MkdirAll(writeDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory recursively %s", uri)
		}

		return os.WriteFile(writeFilePath, bin, f.Mode())
	}

	zipReader, err := zip.NewReader(data, data.Size())
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		data, err := f.Open()
		if err != nil {
			return err
		}
		if err := callbackFunc(f.Name, f.FileInfo(), data); err != nil {
			return err
		}
		_ = data.Close()
	}

	return nil
}
