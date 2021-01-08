// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "time"

type Option func(c *config)

func Title(text string) Option {
	return func(c *config) {
		c.title = text
	}
}

func Size(width, height int) Option {
	return func(c *config) {
		c.width = width
		c.height = height
	}
}

func EnableEvents(events ...Event) Option {
	return func(c *config) {
		for _, e := range events {
			c.eventMask |= e.mask()
		}
	}
}

func DisableCursor() Option {
	return func(c *config) {
		c.cursorDisabled = true
	}
}

func DisableContextMenu() Option {
	return func(c *config) {
		c.contextMenuDisabled = true
	}
}

func FullPage() Option {
	return func(c *config) {
		c.fullPage = true
	}
}

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
