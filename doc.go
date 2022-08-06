// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package canvas communicates via WebSockets with an HTML5 canvas in a web
// browser. The server program sends draw commands to the canvas and the
// canvas sends mouse and keyboard events to the server.
//
// The API doc comments are based on the MDN Web Docs for the [Canvas API]
// by Mozilla Contributors and are licensed under [CC-BY-SA 2.5].
//
// [Canvas API]: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
// [CC-BY-SA 2.5]: https://creativecommons.org/licenses/by-sa/2.5/
package canvas
