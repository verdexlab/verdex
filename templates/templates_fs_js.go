//go:build js
// +build js

package templates_fs

import (
	"embed"
	"io/fs"
)

//go:embed *
var templatesFS embed.FS

func GetTemplatesFs(templatesDirectory string) fs.FS {
	return templatesFS
}
