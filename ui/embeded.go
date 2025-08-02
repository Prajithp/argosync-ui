package ui

import (
	"embed"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:dist
	distFS embed.FS
	//go:embed  dist/index.html
	indexFS    embed.FS
	DistSubFS  = echo.MustSubFS(distFS, "dist")
	IndexSubFS = echo.MustSubFS(indexFS, "dist")
)
