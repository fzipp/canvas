// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Image demonstrates how to draw image data on the canvas.
// It displays a scaling animation of a Go gopher image loaded from PNG data.
//
// Usage:
//
//	image [-http address]
//
// Flags:
//
//	-http  HTTP service address (e.g., '127.0.0.1:8080' or just ':8080').
//	       The default is ':8080'.
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"time"

	"github.com/fzipp/canvas"
)

//go:embed gopher.png
var gopherPNG []byte

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Title:  "ImageData",
		Width:  1280,
		Height: 720,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	img, _, err := image.Decode(bytes.NewBuffer(gopherPNG))
	if err != nil {
		log.Println(err)
		return
	}
	gopher := ctx.CreateImageData(img)
	d := &demo{
		gopher: gopher,
	}
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
		default:
			d.update()
			d.draw(ctx)
			ctx.Flush()
			time.Sleep(20 * time.Millisecond)
		}
	}
}

type demo struct {
	x, y   int
	w, h   int
	gopher *canvas.ImageData
}

func (d *demo) update() {
	d.x = (d.x + 10) % 1280
	d.w = d.x
	d.h = d.x
}

func (d *demo) draw(ctx *canvas.Context) {
	ctx.ClearRect(0, 0, float64(ctx.CanvasWidth()), float64(ctx.CanvasHeight()))
	ctx.DrawImageScaled(d.gopher, float64(d.x), float64(d.y), float64(d.w), float64(d.h))
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
