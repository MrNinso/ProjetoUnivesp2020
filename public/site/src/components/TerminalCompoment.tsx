import React,{ Component } from "react";
import { Terminal} from "xterm"
import { FitAddon }  from "xterm-addon-fit"

import "xterm/css/xterm.css"
import "./TerminalCompoment.css"

function ab2str(buf: ArrayBuffer) {
    return String.fromCharCode.apply(null, Array.from(new Uint8Array(buf)));
}

function CreateTerminal(terminal: HTMLElement | null, id: string) {
    if (terminal != null) {
        let term: Terminal;
        const websocket = new WebSocket(
            "wss://" + window.location.hostname + ":" + window.location.port +
            "/terminal/" + id
        )
        websocket.binaryType = "arraybuffer";

        websocket.onopen = () => {
            term = new Terminal()

            term.loadAddon(FitAddon.prototype)

            term.onData(data => {
                websocket.send(new TextEncoder().encode("\x00" + data));
            })

            term.onResize(evt => {
                websocket.send(new TextEncoder().encode("\x01" + JSON.stringify({cols: evt.cols, rows: evt.rows})))
            })

            term.open(terminal);
            websocket.onmessage = function(evt) {
                if (evt.data instanceof ArrayBuffer) {
                    term.write(ab2str(evt.data));
                } else {
                    console.log(evt.data)
                }
            }

            websocket.onclose = () => {
                term.write("Session terminated");
            }

            websocket.onerror = evt => {
                if (typeof console.log == "function") {
                    console.log(evt)
                }
            }

            term.focus()
            FitAddon.prototype.fit()
        }
    }
}

export class TerminalCompoment extends Component<{ id: string }> {

    render() {
        return <div ref={(t) => CreateTerminal(t, this.props.id)} />
    }
}