const SESSION_SHARE_NOTICE = '// Session is shared with one backend';
const SESSION_SHARE_NOTICE_2 = '// variables persist across snippets even after reload';
const STARTED_LINE = '// Let start with help() command';

export let commandPromptState = $state({
	command: "",
	isEval: false,
	hide: true,
	wrap: false,
	lines: [SESSION_SHARE_NOTICE, SESSION_SHARE_NOTICE_2, STARTED_LINE]
});

interface PrintRequest {
	message: string
}

export const connect = () => {
	// Create a new websocket
	const ws = new WebSocket("/ws");

	ws.addEventListener("print", (message: any) => {
		// Parse the incoming message here
		let data: PrintRequest | null = null;
		try {
			data = JSON.parse(message.data);
		} catch {
			return;
		}
		if (data != null) {
			commandPromptState.lines.push(data.message);
		}
	});

	ws.send("hello")
};
