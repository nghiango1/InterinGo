// Keep it here for later
import {
	commandPromptState,
	WebSocketImpl,
} from "$lib/components/CommandPromptState.svelte";
import type { PageLoad } from "./$types";

export function load({ }: PageLoad) {
	// init commandPromptState new websocket
	if (commandPromptState.ws == null) {
		console.log("[INFO] Create new WS connection to the server");
		commandPromptState.ws = new WebSocketImpl();
	}

	return {
		init: true,
	};
}
