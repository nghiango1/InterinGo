import { createReplRuntime } from "$lib/controller/repl";
import type { CreateReplRuntimeRequest, CreateReplRuntimeResponseSuccess } from "$lib/server/repl";

const SESSION_SHARE_NOTICE = '// Session is shared with one backend';
const SESSION_SHARE_NOTICE_2 = '// variables persist across snippets even after reload';
const STARTED_LINE = '// Let start with help() command';

export interface WebsocketMessageOpen {
	type: "ws_open"
	connId: string
}
export interface WebsocketMessageError {
	type: "ws_error"
	error: string
}
export interface WebsocketMessagePrintEvent {
	type: "ws_print"
	message: string
}

export type WebsocketMessage = WebsocketMessagePrintEvent | WebsocketMessageOpen | WebsocketMessageError

const WS_PATH = '/ws'

export class WebSocketImpl {
	ws: WebSocket
	pingInterval: NodeJS.Timeout | null = null

	constructor() {
		// Create a new websocket
		try {
			this.ws = new WebSocket(WS_PATH);
		} catch {
			const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
			const host = window.location.host; // includes hostname + port

			const ws = `${protocol}//${host}${WS_PATH}`;
			this.ws = new WebSocket(ws);
		}

		// Connection opened
		this.ws.addEventListener("open", (_: Event) => {
			this.ws.send("Hello Server!");
		});

		this.ws.addEventListener('message', (event: MessageEvent) => {
			// Parse the incoming message here
			let data: WebsocketMessage | null = null;
			try {
				data = JSON.parse(event.data);
			} catch {
				return;
			}
			if (data == null) {
				return
			}
			if (data.type == "ws_open") {
				console.log("[INFO] Connection ready: ", data.connId)
				commandPromptState.connId = data.connId

				const req: CreateReplRuntimeRequest = {};
				let runtimeId = window.localStorage.getItem("repl-session")
				if (runtimeId == null) {
					createReplRuntime(req).then(([status, response]) => {
						if (status === 200) {
							commandPromptState.runtimeId = (response as CreateReplRuntimeResponseSuccess).runtimeId;
							window.localStorage.setItem("repl-session", commandPromptState.runtimeId)

							let bindRequest = {
								type: "repl_bind",
								runtimeId: commandPromptState.runtimeId
							}
							this.ws.send(JSON.stringify(bindRequest))

							console.log("[INFO] Connected to a seperated REPL Session", commandPromptState.runtimeId)
						} else {
							console.log("[ERROR] Failed to create seperated REPL Session")
						}
					});
				} else {
					commandPromptState.runtimeId = runtimeId
					let bindRequest = {
						type: "repl_bind",
						runtimeId: commandPromptState.runtimeId
					}
					this.ws.send(JSON.stringify(bindRequest))
				}


			} else if (data.type == "ws_error") {
				console.log("[ERROR] WS send error:", data.error)
			} else if (data.type == "ws_print") {
				commandPromptState.lines.push(data.message);
			}
		});

		this.ws.addEventListener('close', (event: CloseEvent) => {
			console.log("Server connection close!", event.reason)
			if (this.pingInterval) {
				clearInterval(this.pingInterval)
				this.pingInterval = null
			}
		});

		this.ws.addEventListener('error', (event: Event) => {
			console.log("Websocket error", event)
			if (this.pingInterval) {
				clearInterval(this.pingInterval)
				this.pingInterval = null
			}
		});
	}

	// Not expected to be close, just timeout is enough
	close() {
		this.ws.close();
		if (this.pingInterval) {
			clearInterval(this.pingInterval)
			this.pingInterval = null
		}
	}
}

export let commandPromptState = $state({
	ws: null as WebSocketImpl | null,
	connId: null as string | null,
	runtimeId: null as string | null,
	command: "",
	isEval: false,
	hide: true,
	wrap: false,
	lines: [SESSION_SHARE_NOTICE, SESSION_SHARE_NOTICE_2, STARTED_LINE]
});
