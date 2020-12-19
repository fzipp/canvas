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

	const knownImageData = {};
	const offscreenCanvas = {};

    const canvases = document.getElementsByTagName("canvas");
    for (let i = 0; i < canvases.length; i++) {
        const canvas = canvases[i];
        const drawUrl = canvas.dataset.websocketDrawUrl;
        const eventMask = parseInt(canvas.dataset.websocketEventMask, 10);
        if (drawUrl) {
            const absoluteDrawUrl = absoluteWebSocketUrl(drawUrl);
            webSocketCanvas(absoluteDrawUrl, canvas, eventMask);
        }
    }

    function absoluteWebSocketUrl(url) {
        if (url.indexOf("ws://") === 0 || url.indexOf("wss://") === 0) {
            return url;
        }
        const wsUrl = new URL(url, window.location.href);
        wsUrl.protocol = wsUrl.protocol.replace("http", "ws");
        return wsUrl.href;
    }

    function webSocketCanvas(url, canvas, eventMask) {
        const ctx = canvas.getContext("2d");
        const webSocket = new WebSocket(url);
        webSocket.binaryType = "arraybuffer";
        webSocket.onmessage = function (event) {
            const data = event.data;
            let offset = 0;
            const len = data.byteLength;
            while (offset < len) {
                offset += draw(ctx, new DataView(data, offset));
            }
        };
        webSocketCanvasEvents(webSocket, canvas, eventMask);
    }

    function webSocketCanvasEvents(webSocket, canvas, eventMask) {
        if (eventMask & 1) {
            canvas.onmousemove = sendMouseEvent(1);
        }
        if (eventMask & 2) {
            canvas.onmousedown = sendMouseEvent(2);
        }
        if (eventMask & 4) {
            canvas.onmouseup = sendMouseEvent(3);
        }
        if (eventMask & 8) {
            document.onkeypress = sendKeyEvent(4);
        }
        if (eventMask & 16) {
            document.onkeydown = sendKeyEvent(5);
        }
        if (eventMask & 32) {
            document.onkeyup = sendKeyEvent(6);
        }

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
				const buffer = data.buffer.slice(13, 13 + len);
				const array = new Uint8ClampedArray(buffer);
                const imageData = new ImageData(array, width, height);
				knownImageData[id] = imageData;
                const offCanvas = document.createElement("canvas");
                offCanvas.width = width;
                offCanvas.height = height;
                offCanvas.getContext("2d").putImageData(imageData, 0, 0);
			    offscreenCanvas[id] = offCanvas;
				return 13 + len;
			}
			case 13:
                ctx.drawImage(offscreenCanvas[data.getUint32(1)],
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
            case 23:
                ctx.globalAlpha = data.getFloat64(1);
                return 9;
            case 24:
                let mode = null;
                switch (data.getUint8(1)) {
                    case 0:
                        mode = "source-over";
                        break;
                    case 1:
                        mode = "source-in";
                        break;
                    case 2:
                        mode = "source-out";
                        break;
                    case 3:
                        mode = "source-atop";
                        break;
                    case 4:
                        mode = "destination-over";
                        break;
                    case 5:
                        mode = "destination-in";
                        break;
                    case 6:
                        mode = "destination-out";
                        break;
                    case 7:
                        mode = "destination-atop";
                        break;
                    case 8:
                        mode = "lighter";
                        break;
                    case 9:
                        mode = "copy";
                        break;
                    case 10:
                        mode = "xor";
                        break;
                    case 11:
                        mode = "multiply";
                        break;
                    case 12:
                        mode = "screen";
                        break;
                    case 13:
                        mode = "overlay";
                        break;
                    case 14:
                        mode = "darken";
                        break;
                    case 15:
                        mode = "lighten";
                        break;
                    case 16:
                        mode = "color-dodge";
                        break;
                    case 17:
                        mode = "color-burn";
                        break;
                    case 18:
                        mode = "hard-light";
                        break;
                    case 19:
                        mode = "soft-light";
                        break;
                    case 20:
                        mode = "difference";
                        break;
                    case 21:
                        mode = "exclusion";
                        break;
                    case 22:
                        mode = "hue";
                        break;
                    case 23:
                        mode = "saturation";
                        break;
                    case 24:
                        mode = "color";
                        break;
                    case 25:
                        mode = "luminosity";
                        break;
                }
                ctx.globalCompositeOperation = mode;
                return 2;
            case 25:
                ctx.imageSmoothingEnabled = !!data.getUint8(1);
                return 2;
            case 28:
                let cap = null;
                switch (data.getUint8(1)) {
                    case 0:
                        cap = "butt";
                        break;
                    case 1:
                        cap = "round";
                        break;
                    case 2:
                        cap = "square";
                        break;
                }
                ctx.lineCap = cap;
                return 2;
            case 29:
                ctx.lineDashOffset = data.getFloat64(1);
                return 9;
            case 30:
                let join = null;
                switch (data.getUint8(1)) {
                    case 0:
                        join = "miter";
                        break;
                    case 1:
                        join = "round";
                        break;
                    case 2:
                        join = "bevel";
                        break;
                }
                ctx.lineJoin = join;
                return 2;
            case 31:
                ctx.lineTo(data.getFloat64(1), data.getFloat64(9));
                return 17;
            case 32:
                ctx.lineWidth = data.getFloat64(1);
                return 9;
            case 34:
                ctx.miterLimit = data.getFloat64(1);
                return 9;
            case 35:
                ctx.moveTo(data.getFloat64(1), data.getFloat64(9));
                return 17;
			case 36:
                ctx.putImageData(knownImageData[data.getUint32(1)],
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
                let align = null;
                switch (data.getUint8(1)) {
                    case 0:
                        align = "start";
                        break;
                    case 1:
                        align = "end";
                        break;
                    case 2:
                        align = "left";
                        break;
                    case 3:
                        align = "right";
                        break;
                    case 4:
                        align = "center";
                        break;
                }
                ctx.textAlign = align;
                return 2;
            case 54:
                let baseline = null;
                switch (data.getUint8(1)) {
                    case 0:
                        baseline = "alphabetic";
                        break;
                    case 1:
                        baseline = "ideographic";
                        break;
                    case 2:
                        baseline = "top";
                        break;
                    case 3:
                        baseline = "bottom";
                        break;
                    case 4:
                        baseline = "hanging";
                        break;
                    case 5:
                        baseline = "middle";
                        break;
                }
                ctx.textBaseline = baseline;
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
                ctx.putImageData(knownImageData[data.getUint32(1)],
					data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29),
                    data.getFloat64(37), data.getFloat64(45));
				return 53;
			case 63:
                ctx.drawImage(offscreenCanvas[data.getUint32(1)],
					data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29));
				return 37;
			case 64:
                ctx.drawImage(offscreenCanvas[data.getUint32(1)],
					data.getFloat64(5), data.getFloat64(13),
                    data.getFloat64(21), data.getFloat64(29),
                    data.getFloat64(37), data.getFloat64(45),
                    data.getFloat64(53), data.getFloat64(61));
				return 69;
			case 65: {
                const id = data.getUint32(1);
                knownImageData[id] = null;
                offscreenCanvas[id] = null;
				return 5;
            }
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
