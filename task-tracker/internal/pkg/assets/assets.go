package assets

import "embed"

var (
	//go:embed sources
	assets embed.FS
)

func GetAssets() embed.FS {
	return assets
}
