package main

import (
	"embed"
	"io/fs"
)

//go:embed public/*

var StaticFS embed.FS

var FinalStaticFS fs.FS

func init() {
	FinalStaticFS, _ = fs.Sub(StaticFS, "public")
}
