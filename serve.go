// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package canvas communicates via WebSockets with an HTML5 canvas in a web
// browser. The server program sends draw commands to the canvas and the
// canvas sends mouse and keyboard events to the server.
package canvas

import (
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func ListenAndServe(addr string, run func(*Context), options ...Option) error {
	config := configFrom(options)
	http.Handle("/", &htmlHandler{
		config: config,
	})
	http.HandleFunc("/canvas-websocket.js", javaScriptHandler)
	http.Handle("/draw", &drawHandler{
		config: config,
		draw:   run,
	})
	return http.ListenAndServe(addr, nil)
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

const (
	evMouseMove byte = 1 + iota
	evMouseDown
	evMouseUp
	evKeyPress
	evKeyDown
	evKeyUp
	evClick
	evDblClick
	evAuxClick
	evWheel
	evTouchStart
	evTouchMove
	evTouchEnd
	evTouchCancel
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
	case evClick:
		return ClickEvent{decodeMouseEvent(p)}, nil
	case evDblClick:
		return DblClickEvent{decodeMouseEvent(p)}, nil
	case evAuxClick:
		return AuxClickEvent{decodeMouseEvent(p)}, nil
	case evWheel:
		return decodeWheelEvent(p), nil
	case evTouchStart:
		return TouchStartEvent{decodeTouchEvent(p)}, nil
	case evTouchMove:
		return TouchMoveEvent{decodeTouchEvent(p)}, nil
	case evTouchEnd:
		return TouchEndEvent{decodeTouchEvent(p)}, nil
	case evTouchCancel:
		return TouchCancelEvent{decodeTouchEvent(p)}, nil
	}
	return nil, errors.New("unknown event type: '" + string(eventType) + "'")
}

func decodeMouseEvent(p []byte) MouseEvent {
	return MouseEvent{
		Buttons:      MouseButtons(p[1]),
		X:            int(byteOrder.Uint32(p[2:])),
		Y:            int(byteOrder.Uint32(p[6:])),
		modifierKeys: modifierKeys(p[10]),
	}
}

func decodeKeyboardEvent(p []byte) KeyboardEvent {
	keyStringLength := int(byteOrder.Uint32(p[2:]))
	return KeyboardEvent{
		Key:          string(p[6 : 6+keyStringLength]),
		modifierKeys: modifierKeys(p[1]),
	}
}

func decodeWheelEvent(p []byte) WheelEvent {
	return WheelEvent{
		MouseEvent: decodeMouseEvent(p),
		DeltaX:     math.Float64frombits(byteOrder.Uint64(p[11:])),
		DeltaY:     math.Float64frombits(byteOrder.Uint64(p[19:])),
		DeltaZ:     math.Float64frombits(byteOrder.Uint64(p[27:])),
		DeltaMode:  DeltaMode(p[35]),
	}
}

func decodeTouchEvent(p []byte) TouchEvent {
	touches, p := decodeTouchList(p[1:])
	changedTouches, p := decodeTouchList(p)
	targetTouches, p := decodeTouchList(p)
	return TouchEvent{
		Touches:        touches,
		ChangedTouches: changedTouches,
		TargetTouches:  targetTouches,
		modifierKeys:   modifierKeys(p[0]),
	}
}

func decodeTouchList(p []byte) (TouchList, []byte) {
	length := p[0]
	p = p[1:]
	list := make(TouchList, length)
	for i := range list {
		var t Touch
		t, p = decodeTouch(p)
		list[i] = t
	}
	return list, p
}

func decodeTouch(p []byte) (Touch, []byte) {
	return Touch{
		Identifier: byteOrder.Uint32(p),
		X:          int(byteOrder.Uint32(p[4:])),
		Y:          int(byteOrder.Uint32(p[8:])),
	}, p[12:]
}
