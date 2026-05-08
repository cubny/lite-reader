// Package web exposes the embedded static frontend assets.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:public
var assets embed.FS

// FS returns the embedded public/ tree rooted at its top-level entries
// (so a request for "/index.html" maps to "public/index.html" inside the embed).
func FS() fs.FS {
	sub, err := fs.Sub(assets, "public")
	if err != nil {
		panic(err)
	}
	return sub
}
