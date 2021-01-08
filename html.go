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
    <style>
      * {
        margin: 0;
        padding: 0;
      }
      body, html {
        height: 100%;
      }
      .full-page {
        position: absolute;
        width: 100%;
        height: 100%;
      }
    </style>
  </head>
  <body>
    <noscript><p>Please enable JavaScript in your browser.</p></noscript>
    <canvas width="{{.Width}}" height="{{.Height}}"
            style="cursor: {{if .CursorDisabled}}none{{else}}default{{end}}"
            class="{{if .FullPage}}full-page{{end}}"
            data-websocket-draw-url="{{.DrawURL}}"
            data-websocket-event-mask="{{.EventMask}}"
            data-websocket-reconnect-interval="{{.ReconnectInterval}}"
            data-disable-context-menu="{{.ContextMenuDisabled}}"></canvas>
  </body>
</html>
`))
}
