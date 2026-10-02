package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
)

// icon.png is rendered from icon.svg (scripts/render-icons.sh)
//
//go:embed icon.png
var iconPNG []byte

func appIcon() image.Image {
	img, err := png.Decode(bytes.NewReader(iconPNG))
	if err != nil {
		return nil
	}
	return img
}
