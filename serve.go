// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	// Package embed is used to embed the HTML template and JavaScript files.
	_ "embed"
	"fmt"
	"html/template"
	"image/color"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	//go:embed web/canvas-websocket.js
	javaScriptCode []byte

	//go:embed web/index.html.tmpl
	indexHTMLCode     string
	indexHTMLTemplate = template.Must(template.New("index.html.tmpl").Parse(indexHTMLCode))
)

// ListenAndServe listens on the TCP network address addr and serves
// an HTML page on "/" with a canvas that connects to the server via
// WebSockets on a "/draw" endpoint. It also serves a JavaScript file on
// "/canvas-websocket.js" that is used by the HTML page to provide this
// functionality.
//
// The run function is called when a client canvas connects to the server.
// The Context parameter of the run function allows the server to send draw
// commands to the canvas and receive events from the canvas. Each instance
// of the run function runs on its own goroutine, so the run function should
// not access shared state without proper synchronization.
//
// The options configure various aspects the canvas such as its size, which
// events to handle etc.
func ListenAndServe(addr string, run func(*Context), options ...Option) error {
	return http.ListenAndServe(addr, NewServeMux(run, options...))
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS / WSS connections. Additionally, files containing a
// certificate and matching private key for the server must be provided. If the
// certificate is signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and the CA's
// certificate.
func ListenAndServeTLS(addr, certFile, keyFile string, run func(*Context), options ...Option) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, NewServeMux(run, options...))
}

// NewServeMux creates a http.ServeMux as used by ListenAndServe.
func NewServeMux(run func(*Context), options ...Option) *http.ServeMux {
	config := configFrom(options)
	mux := http.NewServeMux()
	mux.Handle("/", &htmlHandler{
		config: config,
	})
	mux.HandleFunc("/canvas-websocket.js", javaScriptHandler)
	mux.Handle("/draw", &drawHandler{
		config: config,
		draw:   run,
	})
	return mux
}

type htmlHandler struct {
	config config
}

func (h *htmlHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	model := map[string]interface{}{
		"DrawURL":             template.URL("draw"),
		"Width":               h.config.width,
		"Height":              h.config.height,
		"Title":               h.config.title,
		"BackgroundColor":     template.CSS(rgbaString(h.config.backgroundColor)),
		"EventMask":           h.config.eventMask,
		"CursorDisabled":      h.config.cursorDisabled,
		"ContextMenuDisabled": h.config.contextMenuDisabled,
		"FullPage":            h.config.fullPage,
		"ReconnectInterval":   int64(h.config.reconnectInterval / time.Millisecond),
	}
	err := indexHTMLTemplate.Execute(w, model)
	if err != nil {
		log.Println(err)
		return
	}
}

func rgbaString(c color.Color) string {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("rgba(%d, %d, %d, %g)", clr.R, clr.G, clr.B, float64(clr.A)/255)
}

func javaScriptHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	_, err := w.Write(javaScriptCode)
	if err != nil {
		log.Println(err)
		return
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type drawHandler struct {
	config config
	draw   func(*Context)
}

func (h *drawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	events := make(chan Event)
	defer close(events)
	draws := make(chan []byte)
	defer close(draws)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go readMessages(conn, events, &wg)
	go writeMessages(conn, draws, &wg)

	ctx := newContext(draws, events, h.config)
	go func() {
		defer wg.Done()
		h.draw(ctx)
	}()

	wg.Wait()
	wg.Add(1)
	events <- CloseEvent{}
	wg.Wait()
}

func writeMessages(conn *websocket.Conn, messages <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message := <-messages
		err := conn.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			break
		}
	}
}

func readMessages(conn *websocket.Conn, events chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.BinaryMessage {
			continue
		}
		event, err := decodeEvent(p)
		if err != nil {
			continue
		}
		events <- event
	}
}
