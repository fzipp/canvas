// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"time"
)

// Option represents an option that configures the canvas.
// The functions ListenAndServe, ListenAndServeTLS and NewServeMux take
// a variable number of options. Options are created by option functions
// like Title, Size, EnableEvents, etc.
type Option func(c *config)

// Title returns an option that sets the title of the browser tab/window.
func Title(text string) Option {
	return func(c *config) {
		c.title = text
	}
}

// Size returns an option that sets the width and height of the canvas.
func Size(width, height int) Option {
	return func(c *config) {
		c.width = width
		c.height = height
	}
}

// EnableEvents returns an option that enables transmission of the given event
// types from the client to the server.
func EnableEvents(events ...Event) Option {
	return func(c *config) {
		for _, e := range events {
			c.eventMask |= e.mask()
		}
	}
}

// DisableCursor returns an option that hides the mouse cursor on the canvas.
func DisableCursor() Option {
	return func(c *config) {
		c.cursorDisabled = true
	}
}

// DisableContextMenu returns an option that disables the context menu on
// the canvas.
func DisableContextMenu() Option {
	return func(c *config) {
		c.contextMenuDisabled = true
	}
}

// ScaleFullPage returns an option that scales the canvas to the full extent of
// the page (horizontally, vertically, or both) in the browser window.
func ScaleFullPage(fullWidth, fullHeight bool) Option {
	return func(c *config) {
		c.fullPageWidth = fullWidth
		c.fullPageHeight = fullHeight
	}
}

// Reconnect returns an option that configures the client to reconnect after
// the given duration if the websocket connection was lost.
// The client tries to reconnect repeatedly until it is successful.
func Reconnect(interval time.Duration) Option {
	return func(c *config) {
		c.reconnectInterval = interval
	}
}

// BackgroundColor returns an option that configures the background color of
// the served HTML page.
func BackgroundColor(c color.Color) Option {
	return func(cfg *config) {
		cfg.backgroundColor = c
	}
}

type config struct {
	title               string
	width               int
	height              int
	backgroundColor     color.Color
	eventMask           eventMask
	cursorDisabled      bool
	contextMenuDisabled bool
	fullPageWidth       bool
	fullPageHeight      bool
	reconnectInterval   time.Duration
}

func configFrom(options []Option) config {
	c := &config{
		width: 300, height: 150,
		backgroundColor: color.White,
	}
	for _, opt := range options {
		opt(c)
	}
	return *c
}
