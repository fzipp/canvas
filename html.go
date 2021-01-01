// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "html/template"

var htmlTemplate *template.Template

func init() {
	htmlTemplate = template.Must(template.New("index.html").Parse(`
<!doctype html>
<html>
  <head>
    <title>{{.Title}}</title>
    <script src="canvas-websocket.js"></script>
    <style>body {margin: 0}</style>
  </head>
  <body>
    <canvas width="{{.Width}}" height="{{.Height}}"
            style="cursor: {{if .CursorDisabled}}none{{else}}default{{end}}"
            data-websocket-draw-url="{{.DrawURL}}"
            data-websocket-event-mask="{{.EventMask}}"
            data-disable-context-menu="{{.ContextMenuDisabled}}"></canvas>
  </body>
</html>
`))
}
