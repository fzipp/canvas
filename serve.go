// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package canvas communicates via WebSockets with an HTML5 canvas in a web
// browser. The server program sends draw commands to the canvas and the
// canvas sends mouse and keyboard events to the server.
package canvas

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func ListenAndServe(addr string, run func(*Context), options ...Option) error {
	return http.ListenAndServe(addr, NewServeMux(run, options...))
}

func ListenAndServeTLS(addr, certFile, keyFile string, run func(*Context), options ...Option) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, NewServeMux(run, options...))
}

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

func (h *htmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	model := map[string]interface{}{
		"DrawURL":             template.URL("draw"),
		"Width":               h.config.width,
		"Height":              h.config.height,
		"Title":               h.config.title,
		"EventMask":           h.config.eventMask,
		"CursorDisabled":      h.config.cursorDisabled,
		"ContextMenuDisabled": h.config.contextMenuDisabled,
		"FullPage":            h.config.fullPage,
		"ReconnectInterval":   int64(h.config.reconnectInterval / time.Millisecond),
	}
	err := htmlTemplate.Execute(w, model)
	if err != nil {
		log.Println(err)
		return
	}
}

func javaScriptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	err := javaScriptTemplate.Execute(w, nil)
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
	quit := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(2)
	go readMessages(conn, events, &wg)
	go writeMessages(conn, draws, &wg)

	ctx := newContext(draws, events, quit, h.config)
	go func() {
		defer wg.Done()
		h.draw(ctx)
	}()

	wg.Wait()
	wg.Add(1)
	close(quit)
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
