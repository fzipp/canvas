// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example was ported from the MDN canvas tutorial [An animated clock].
// Original copyright: Any copyright is dedicated to the [Public Domain].
//
// [An animated clock]: https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Basic_animations#An_animated_clock
// [Public Domain]: https://creativecommons.org/publicdomain/zero/1.0/

// Clock draws an animated clock, showing the current time.
//
// Usage:
//
//	clock [-http address]
//
// Flags:
//
//	-http  HTTP service address (e.g., '127.0.0.1:8080' or just ':8080').
//	       The default is ':8080'.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Title:  "Clock",
		Width:  150,
		Height: 150,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
		default:
			drawClock(ctx)
			ctx.Flush()
			time.Sleep(time.Second / 2)
		}
	}
}

func drawClock(ctx *canvas.Context) {
	now := time.Now()

	ctx.Save()
	ctx.ClearRect(0, 0, 150, 150)
	ctx.Translate(75, 75)
	ctx.Scale(0.4, 0.4)
	ctx.Rotate(-math.Pi / 2)
	ctx.SetStrokeStyle(color.Black)
	ctx.SetFillStyle(color.White)
	ctx.SetLineWidth(8)
	ctx.SetLineCap(canvas.CapRound)

	// Hour marks
	ctx.Save()
	for i := 0; i < 12; i++ {
		ctx.BeginPath()
		ctx.Rotate(math.Pi / 6)
		ctx.MoveTo(100, 0)
		ctx.LineTo(120, 0)
		ctx.Stroke()
	}
	ctx.Restore()

	// Minute marks
	ctx.Save()
	ctx.SetLineWidth(5)
	for i := 0; i < 60; i++ {
		if i%5 != 0 {
			ctx.BeginPath()
			ctx.MoveTo(117, 0)
			ctx.LineTo(120, 0)
			ctx.Stroke()
		}
		ctx.Rotate(math.Pi / 30)
	}
	ctx.Restore()

	second := float64(now.Second())
	minute := float64(now.Minute())
	hour := float64(now.Hour())
	if hour >= 12 {
		hour = hour - 12
	}

	ctx.SetFillStyle(color.Black)

	// write Hours
	ctx.Save()
	ctx.Rotate(hour*(math.Pi/6) + (math.Pi/360)*minute + (math.Pi/21600)*second)
	ctx.SetLineWidth(14)
	ctx.BeginPath()
	ctx.MoveTo(-20, 0)
	ctx.LineTo(80, 0)
	ctx.Stroke()
	ctx.Restore()

	// write Minutes
	ctx.Save()
	ctx.Rotate((math.Pi/30)*minute + (math.Pi/1800)*second)
	ctx.SetLineWidth(10)
	ctx.BeginPath()
	ctx.MoveTo(-28, 0)
	ctx.LineTo(112, 0)
	ctx.Stroke()
	ctx.Restore()

	// Write seconds
	ctx.Save()
	ctx.Rotate(second * math.Pi / 30)
	ctx.SetStrokeStyleString("#D40000")
	ctx.SetFillStyleString("#D40000")
	ctx.SetLineWidth(6)
	ctx.BeginPath()
	ctx.MoveTo(-30, 0)
	ctx.LineTo(83, 0)
	ctx.Stroke()
	ctx.BeginPath()
	ctx.Arc(0, 0, 10, 0, math.Pi*2, true)
	ctx.Fill()
	ctx.BeginPath()
	ctx.Arc(95, 0, 10, 0, math.Pi*2, true)
	ctx.Stroke()
	ctx.SetFillStyle(color.Transparent)
	ctx.Arc(0, 0, 3, 0, math.Pi*2, true)
	ctx.Fill()
	ctx.Restore()

	ctx.BeginPath()
	ctx.SetLineWidth(14)
	ctx.SetStrokeStyleString("#325FA2")
	ctx.Arc(0, 0, 142, 0, math.Pi*2, true)
	ctx.Stroke()

	ctx.Restore()
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
