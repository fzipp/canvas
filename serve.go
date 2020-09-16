// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func ListenAndServe(addr string, width, height int, title string, run func(*Context), eventMask SendEventMask) {
	http.Handle("/", &htmlHandler{
		width: width, height: height,
		title: title, eventMask: eventMask,
	})
	http.HandleFunc("/canvas-websocket.js", javaScriptHandler)
	http.Handle("/run", &drawHandler{
		width: width, height: height, draw: run,
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println(err)
	}
}

type htmlHandler struct {
	width     int
	height    int
	title     string
	eventMask SendEventMask
}

func (h *htmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	model := map[string]interface{}{
		"Width":         h.width,
		"Height":        h.height,
		"Title":         h.title,
		"DrawURL":       template.URL("ws://localhost:8080/draw"),
		"SendEventMask": h.eventMask,
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
	width  int
	height int
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

	ctx := newContext(h.width, h.height, draws, events, quit)
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

const (
	evMouseMove byte = 1 + iota
	evMouseDown
	evMouseUp
	evKeyPress
	evKeyDown
	evKeyUp
)

func decodeEvent(p []byte) (Event, error) {
	eventType := p[0]
	switch eventType {
	case evMouseMove:
		return MouseMoveEvent{decodeMouseEvent(p)}, nil
	case evMouseDown:
		return MouseDownEvent{decodeMouseEvent(p)}, nil
	case evMouseUp:
		return MouseUpEvent{decodeMouseEvent(p)}, nil
	case evKeyPress:
		return KeyPressEvent{decodeKeyboardEvent(p)}, nil
	case evKeyDown:
		return KeyDownEvent{decodeKeyboardEvent(p)}, nil
	case evKeyUp:
		return KeyUpEvent{decodeKeyboardEvent(p)}, nil
	}
	return nil, errors.New("unknown event type: '" + string(eventType) + "'")
}

func decodeMouseEvent(p []byte) MouseEvent {
	return MouseEvent{
		Button:  int(p[1]),
		X:       int(byteOrder.Uint32(p[2:])),
		Y:       int(byteOrder.Uint32(p[6:])),
		modKeys: modifierKey(p[10]),
	}
}

func decodeKeyboardEvent(p []byte) KeyboardEvent {
	keyStringLength := int(byteOrder.Uint32(p[2:]))
	return KeyboardEvent{
		Key:     string(p[6 : 6+keyStringLength]),
		modKeys: modifierKey(p[1]),
	}
}
