// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "time"

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

// FullPage returns an option that resizes the canvas to the full extent of the
// page in the browser window.
func FullPage() Option {
	return func(c *config) {
		c.fullPage = true
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

type config struct {
	title               string
	width               int
	height              int
	eventMask           eventMask
	cursorDisabled      bool
	contextMenuDisabled bool
	fullPage            bool
	reconnectInterval   time.Duration
}

func configFrom(options []Option) config {
	c := &config{}
	for _, opt := range options {
		opt(c)
	}
	return *c
}
