//go:build !js
// +build !js

package templates_fs

import (
	"io/fs"
	"os"
)

func GetTemplatesFs(templatesDirectory string) fs.FS {
	return os.DirFS(templatesDirectory)
}
