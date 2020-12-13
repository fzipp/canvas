// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

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

func EnableEvents(eventTypes ...SendEventMask) Option {
	return func(c *config) {
		for _, eventType := range eventTypes {
			c.eventMask |= eventType
		}
	}
}

func DisableCursor() Option {
	return func(c *config) {
		c.cursorDisabled = true
	}
}

type config struct {
	title          string
	width          int
	height         int
	eventMask      SendEventMask
	cursorDisabled bool
}

func configFrom(options []Option) config {
	c := &config{}
	for _, opt := range options {
		opt(c)
	}
	return *c
}
