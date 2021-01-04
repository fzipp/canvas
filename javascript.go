// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "text/template"

var javaScriptTemplate *template.Template

func init() {
	javaScriptTemplate = template.Must(template.New("canvas-websocket.js").Parse(`
document.addEventListener("DOMContentLoaded", function () {
    "use strict";

    const allocImageData = {};
    const allocOffscreenCanvas = {};
    const allocGradient = {};
    const allocPattern = {};

    const enumRepetition = ["repeat", "repeat-x", "repeat-y", "no-repeat"];

    const enumCompositeOperation = [
        "source-over", "source-in", "source-out", "source-atop",
        "destination-over", "destination-in", "destination-out",
        "destination-atop", "lighter", "copy", "xor", "multiply", "screen",
        "overlay", "darken", "lighten", "color-dodge", "color-burn",
        "hard-light", "soft-light", "difference", "exclusion", "hue",
        "saturation", "color", "luminosity"
    ];

    const enumLineCap = ["butt", "round", "square"];

    const enumLineJoin = ["miter", "round", "bevel"];

    const enumTextAlign = ["start", "end", "left", "right", "center"];

    const enumTextBaseline = [
        "alphabetic", "ideographic", "top", "bottom", "middle"
    ];

    const canvases = document.getElementsByTagName("canvas");
    for (let i = 0; i < canvases.length; i++) {
        const canvas = canvases[i];
        const config = configFrom(canvas.dataset);
        if (config.drawUrl) {
            webSocketCanvas(canvas, config);
            if (config.contextMenuDisabled) {
                disableContextMenu(canvas);
            }
        }
    }

    function configFrom(dataset) {
        return {
            drawUrl: absoluteWebSocketUrl(dataset["websocketDrawUrl"]),
            eventMask: parseInt(dataset["websocketEventMask"], 10) || 0,
            reconnectInterval: parseInt(dataset["websocketReconnectInterval"], 10) || 0,
            contextMenuDisabled: (dataset["disableContextMenu"] === "true")
        };
    }

    function absoluteWebSocketUrl(url) {
        if (!url) {
            return null;
        }
        if (url.indexOf("ws://") === 0 || url.indexOf("wss://") === 0) {
            return url;
        }
        const wsUrl = new URL(url, window.location.href);
        wsUrl.protocol = wsUrl.protocol.replace("http", "ws");
        return wsUrl.href;
    }

    function webSocketCanvas(canvas, config) {
        const ctx = canvas.getContext("2d");
        const webSocket = new WebSocket(config.drawUrl);
        let handlers = {};
        webSocket.binaryType = "arraybuffer";
        webSocket.addEventListener("open", function () {
            handlers = addEventListeners(canvas, config.eventMask, webSocket);
        });
        webSocket.addEventListener("error", function () {
            webSocket.close();
        });
        webSocket.addEventListener("close", function () {
            removeEventListeners(canvas, handlers);
            if (!config.reconnectInterval) {
                return;
            }
            setTimeout(function () {
                webSocketCanvas(canvas, config);
            }, config.reconnectInterval);
        });
        webSocket.addEventListener("message", function (event) {
            const data = event.data;
            let offset = 0;
            const len = data.byteLength;
            while (offset < len) {
                offset += draw(ctx, new DataView(data, offset));
            }
        });
    }

    function addEventListeners(canvas, eventMask, webSocket) {
        const handlers = {};

        if (eventMask & 1) {
            handlers["mousemove"] = sendMouseEvent(1);
        }
        if (eventMask & 2) {
            handlers["mousedown"] = sendMouseEvent(2);
        }
        if (eventMask & 4) {
            handlers["onmouseup"] = sendMouseEvent(3);
        }
        if (eventMask & 8) {
            handlers["keypress"] = sendKeyEvent(4);
        }
        if (eventMask & 16) {
            handlers["keydown"] = sendKeyEvent(5);
        }
        if (eventMask & 32) {
            handlers["keyup"] = sendKeyEvent(6);
        }
        if (eventMask & 64) {
            handlers["click"] = sendMouseEvent(7);
        }
        if (eventMask & 128) {
            handlers["dblclick"] = sendMouseEvent(8);
        }
        if (eventMask & 256) {
            handlers["auxclick"] = sendMouseEvent(9);
        }

        Object.keys(handlers).forEach(function (type) {
            const target = (type.indexOf("key") !== 0) ? canvas : document;
            target.addEventListener(type, handlers[type]);
        });

        const rect = canvas.getBoundingClientRect();

        function sendMouseEvent(eventType) {
            return function (event) {
                const eventMessage = new ArrayBuffer(11);
                const data = new DataView(eventMessage);
                data.setUint8(0, eventType);
                data.setUint8(1, event.buttons);
                data.setUint32(2, event.clientX - rect.left);
                data.setUint32(6, event.clientY - rect.top);
                data.setUint8(10, encodeModifierKeys(event));
                webSocket.send(eventMessage);
            };
        }

        function sendKeyEvent(eventType) {
            return function (event) {
                const keyBytes = new TextEncoder().encode(event.key);
                const eventMessage = new ArrayBuffer(6 + keyBytes.byteLength);
                const data = new DataView(eventMessage);
                data.setUint8(0, eventType);
                data.setUint8(1, encodeModifierKeys(event));
                data.setUint32(2, keyBytes.byteLength);
                for (let i = 0; i < keyBytes.length; i++) {
                    data.setUint8(6 + i, keyBytes[i]);
                }
                webSocket.send(eventMessage);
            };
        }

        return handlers;
    }

    function removeEventListeners(canvas, handlers) {
        Object.keys(handlers).forEach(function (type) {
            const target = (type.indexOf("key") !== 0) ? canvas : document;
            target.removeEventListener(type, handlers[type]);
        });
    }

    function disableContextMenu(canvas) {
        canvas.addEventListener("contextmenu", function (e) {
            e.preventDefault();
        }, false);
    }

    function encodeModifierKeys(event) {
        let modifiers = 0;
        if (event.altKey) {
            modifiers |= 1;
        }
        if (event.shiftKey) {
            modifiers |= 2;
        }
        if (event.ctrlKey) {
            modifiers |= 4;
        }
        if (event.metaKey) {
            modifiers |= 8;
        }
        return modifiers;
    }

    function draw(ctx, data) {
        switch (data.getUint8(0)) {
            case 1:
                ctx.arc(
                    data.getFloat64(1), data.getFloat64(9), data.getFloat64(17),
                    data.getFloat64(25), data.getFloat64(33), !!data.getUint8(41));
                return 42;
            case 2:
                ctx.arcTo(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25),
                    data.getFloat64(33));
                return 41;
            case 3:
                ctx.beginPath();
                return 1;
            case 4:
                ctx.bezierCurveTo(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25),
                    data.getFloat64(33), data.getFloat64(41));
                return 49;
            case 5:
                ctx.clearRect(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25));
                return 33;
            case 6:
                ctx.clip();
                return 1;
            case 8: {
                const id = data.getUint32(1);
                const width = data.getUint32(5);
                const height = data.getUint32(9);
                const len = width * height * 4;
                const bufferOffset = data.byteOffset + 13;
                const buffer = data.buffer.slice(bufferOffset, bufferOffset + len);
                const array = new Uint8ClampedArray(buffer);
                const imageData = new ImageData(array, width, height);
                allocImageData[id] = imageData;
                const offCanvas = document.createElement("canvas");
                offCanvas.width = width;
                offCanvas.height = height;
                offCanvas.getContext("2d").putImageData(imageData, 0, 0);
                allocOffscreenCanvas[id] = offCanvas;
                return 13 + len;
            }
            case 9: {
                const id = data.getUint32(1);
                const x0 = data.getFloat64(5);
                const y0 = data.getFloat64(13);
                const x1 = data.getFloat64(21);
                const y1 = data.getFloat64(29);
                allocGradient[id] = ctx.createLinearGradient(x0, y0, x1, y1);
                return 37;
            }
            case 10: {
                const id = data.getUint32(1);
                const image = allocOffscreenCanvas[data.getUint32(5)];
                const repetition = enumRepetition[data.getUint8(9)];
                allocPattern[id] = ctx.createPattern(image, repetition);
                return 10;
            }
            case 11: {
                const id = data.getUint32(1);
                const x0 = data.getFloat64(5);
                const y0 = data.getFloat64(13);
                const r0 = data.getFloat64(21);
                const x1 = data.getFloat64(29);
                const y1 = data.getFloat64(37);
                const r1 = data.getFloat64(45);
                allocGradient[id] = ctx.createRadialGradient(x0, y0, r0, x1, y1, r1);
                return 53;
            }
            case 13:
                ctx.drawImage(allocOffscreenCanvas[data.getUint32(1)],
                    data.getFloat64(5), data.getFloat64(13));
                return 21;
            case 14:
                ctx.ellipse(
                    data.getFloat64(1), data.getFloat64(9), data.getFloat64(17),
                    data.getFloat64(25), data.getFloat64(33), data.getFloat64(41),
                    data.getFloat64(49), !!data.getUint8(57));
                return 58;
            case 15:
                ctx.fill();
                return 1;
            case 16:
                ctx.fillRect(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25));
                return 33;
            case 17: {
                ctx.fillStyle = getRGBA(data, 1);
                return 5;
            }
            case 18: {
                const text = getString(data, 17);
                ctx.fillText(text.value, data.getFloat64(1), data.getFloat64(9));
                return 17 + text.byteLen;
            }
            case 19: {
                const font = getString(data, 1);
                ctx.font = font.value;
                return 1 + font.byteLen;
            }
            case 20: {
                const id = data.getUint32(1);
                const gradient = allocGradient[id];
                gradient.addColorStop(data.getFloat64(5), getRGBA(data, 13));
                return 17;
            }
            case 21: {
                const id = data.getUint32(1);
                const offset = data.getFloat64(5);
                const color = getString(data, 13);
                const gradient = allocGradient[id];
                gradient.addColorStop(offset, color.value);
                return 13 + color.byteLen;
            }
            case 22:
                ctx.fillStyle = allocGradient[data.getUint32(1)];
                return 5;
            case 23:
                ctx.globalAlpha = data.getFloat64(1);
                return 9;
            case 24:
                ctx.globalCompositeOperation = enumCompositeOperation[data.getUint8(1)];
                return 2;
            case 25:
                ctx.imageSmoothingEnabled = !!data.getUint8(1);
                return 2;
            case 26:
                ctx.strokeStyle = allocGradient[data.getUint32(1)];
                return 5;
            case 27: {
                const id = data.getUint32(1);
                allocPattern[id] = null;
                return 5;
            }
            case 28:
                ctx.lineCap = enumLineCap[data.getUint8(1)];
                return 2;
            case 29:
                ctx.lineDashOffset = data.getFloat64(1);
                return 9;
            case 30:
                ctx.lineJoin = enumLineJoin[data.getUint8(1)];
                return 2;
            case 31:
                ctx.lineTo(data.getFloat64(1), data.getFloat64(9));
                return 17;
            case 32:
                ctx.lineWidth = data.getFloat64(1);
                return 9;
            case 33: {
                const id = data.getUint32(1);
                allocGradient[id] = null;
                return 5;
            }
            case 34:
                ctx.miterLimit = data.getFloat64(1);
                return 9;
            case 35:
                ctx.moveTo(data.getFloat64(1), data.getFloat64(9));
                return 17;
            case 36:
                ctx.putImageData(allocImageData[data.getUint32(1)],
                    data.getFloat64(5), data.getFloat64(13));
                return 21;
            case 37:
                ctx.quadraticCurveTo(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25));
                return 33;
            case 38:
                ctx.rect(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25));
                return 33;
            case 39:
                ctx.restore();
                return 1;
            case 40:
                ctx.rotate(data.getFloat64(1));
                return 9;
            case 41:
                ctx.save();
                return 1;
            case 42:
                ctx.scale(data.getFloat64(1), data.getFloat64(9));
                return 17;
            case 43: {
                const segments = [];
                const len = data.getUint32(1);
                for (let i = 0; i < len; i++) {
                    segments.push(data.getFloat64(5 + i * 8));
                }
                ctx.setLineDash(segments);
                return 5 + (len * 8);
            }
            case 44:
                ctx.setTransform(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25),
                    data.getFloat64(33), data.getFloat64(41));
                return 49;
            case 45:
                ctx.shadowBlur = data.getFloat64(1);
                return 9;
            case 46: {
                ctx.shadowColor = getRGBA(data, 1);
                return 5;
            }
            case 47:
                ctx.shadowOffsetX = data.getFloat64(1);
                return 9;
            case 48:
                ctx.shadowOffsetY = data.getFloat64(1);
                return 9;
            case 49:
                ctx.stroke();
                return 1;
            case 50:
                ctx.strokeRect(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25));
                return 33;
            case 51: {
                ctx.strokeStyle = getRGBA(data, 1);
                return 5;
            }
            case 52: {
                const text = getString(data, 17);
                ctx.strokeText(text.value, data.getFloat64(1), data.getFloat64(9));
                return 17 + text.byteLen;
            }
            case 53:
                ctx.textAlign = enumTextAlign[data.getUint8(1)];
                return 2;
            case 54:
                ctx.textBaseline = enumTextBaseline[data.getUint8(1)];
                return 2;
            case 55:
                ctx.transform(
                    data.getFloat64(1), data.getFloat64(9),
                    data.getFloat64(17), data.getFloat64(25),
                    data.getFloat64(33), data.getFloat64(41));
                return 49;
            case 56:
                ctx.translate(data.getFloat64(1), data.getFloat64(9));
                return 17;
            case 57: {
                const text = getString(data, 25);
                ctx.fillText(text.value, data.getFloat64(1), data.getFloat64(9), data.getFloat64(17));
                return 25 + text.byteLen;
            }
            case 58: {
                const text = getString(data, 25);
                ctx.strokeText(text.value, data.getFloat64(1), data.getFloat64(9), data.getFloat64(17));
                return 25 + text.byteLen;
            }
            case 59: {
                const color = getString(data, 1);
                ctx.fillStyle = color.value;
                return 1 + color.byteLen;
            }
            case 60: {
                const color = getString(data, 1);
                ctx.strokeStyle = color.value;
                return 1 + color.byteLen;
            }
            case 61: {
                const color = getString(data, 1);
                ctx.shadowColor = color.value;
                return 1 + color.byteLen;
            }
            case 62:
                ctx.putImageData(allocImageData[data.getUint32(1)],
                    data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29),
                    data.getFloat64(37), data.getFloat64(45));
                return 53;
            case 63:
                ctx.drawImage(allocOffscreenCanvas[data.getUint32(1)],
                    data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29));
                return 37;
            case 64:
                ctx.drawImage(allocOffscreenCanvas[data.getUint32(1)],
                    data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29),
                    data.getFloat64(37), data.getFloat64(45),
                    data.getFloat64(53), data.getFloat64(61));
                return 69;
            case 65: {
                const id = data.getUint32(1);
                allocImageData[id] = null;
                allocOffscreenCanvas[id] = null;
                return 5;
            }
            case 66:
                ctx.fillStyle = allocPattern[data.getUint32(1)];
                return 5;
            case 67:
                ctx.strokeStyle = allocPattern[data.getUint32(1)];
                return 5;
        }
        return 1;
    }

    function getString(data, offset) {
        const stringLen = data.getUint32(offset);
        const stringBegin = data.byteOffset + offset + 4;
        const stringEnd = stringBegin + stringLen;
        return {
            value: new TextDecoder().decode(data.buffer.slice(stringBegin, stringEnd)),
            byteLen: 4 + stringLen
        };
    }

    function getRGBA(data, offset) {
        return "rgba(" +
            data.getUint8(offset) + ", " +
            data.getUint8(offset + 1) + ", " +
            data.getUint8(offset + 2) + ", " +
            data.getUint8(offset + 3) / 255 + ")";
    }
});
`))
}
