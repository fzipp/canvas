// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example was ported from the MDN canvas tutorial: "An animated clock"
// https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Basic_animations#An_animated_clock
// Original copyright:
//
// Any copyright is dedicated to the Public Domain.
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	port := ":8080"
	fmt.Println("Listening on http://localhost" + port)
	err := canvas.ListenAndServe(port, run,
		canvas.Size(800, 600),
		canvas.Title("Clock"),
	)
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

	sec := float64(now.Second())
	min := float64(now.Minute())
	hr := float64(now.Hour())
	if hr >= 12 {
		hr = hr - 12
	}

	ctx.SetFillStyle(color.Black)

	// write Hours
	ctx.Save()
	ctx.Rotate(hr*(math.Pi/6) + (math.Pi/360)*min + (math.Pi/21600)*sec)
	ctx.SetLineWidth(14)
	ctx.BeginPath()
	ctx.MoveTo(-20, 0)
	ctx.LineTo(80, 0)
	ctx.Stroke()
	ctx.Restore()

	// write Minutes
	ctx.Save()
	ctx.Rotate((math.Pi/30)*min + (math.Pi/1800)*sec)
	ctx.SetLineWidth(10)
	ctx.BeginPath()
	ctx.MoveTo(-28, 0)
	ctx.LineTo(112, 0)
	ctx.Stroke()
	ctx.Restore()

	// Write seconds
	ctx.Save()
	ctx.Rotate(sec * math.Pi / 30)
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
