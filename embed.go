package zinc

import (
	"embed"
	"io/fs"
)

//go:embed web/dist
var embedFrontend embed.FS

func GetFrontendAssets() (fs.FS, error) {
	f, err := fs.Sub(embedFrontend, "web/dist")
	if err != nil {
		return nil, err
	}

	return f, nil
}
