package web

import "embed"

//go:embed pages/*.html
var Templates embed.FS
