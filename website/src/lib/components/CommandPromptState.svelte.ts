const SESSION_SHARE_NOTICE = '// Session is shared with one backend';
const SESSION_SHARE_NOTICE_2 = '// variables persist across snippets even after reload';
const STARTED_LINE = '// Let start with help() command';

export interface PrintRequest {
	message: string
}

class WebSocketImpl {
	ws: WebSocket
	pingInterval: NodeJS.Timeout | null = null

	constructor() {
		// Create a new websocket
		this.ws = new WebSocket('/ws');

		// Connection opened
		this.ws.addEventListener("open", (_: Event) => {
			this.ws.send("Hello Server!");

			// Keep connection alive, else server will clean up this session
			// ping every 1 mins
			this.pingInterval = setInterval(() => {
				console.log(`SENT: ping`);
				this.ws.send("ping");
			}, 1000 * 60);
		});

		this.ws.addEventListener('message', (event: MessageEvent) => {
			// Parse the incoming message here
			let data: PrintRequest | null = null;
			try {
				data = JSON.parse(event.data);
			} catch {
				return;
			}
			if (data != null) {
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
	ws: new WebSocketImpl(),
	command: "",
	isEval: false,
	hide: true,
	wrap: false,
	lines: [SESSION_SHARE_NOTICE, SESSION_SHARE_NOTICE_2, STARTED_LINE]
});
