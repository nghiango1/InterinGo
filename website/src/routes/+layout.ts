export const prerender = true;
export const ssr = false;

// Keep it here for later
import type { LayoutLoad } from "./$types";
import { commandPromptState, WebSocketImpl } from "$lib/components/CommandPromptState.svelte";

export function load({ }: LayoutLoad) {
	console.log("[INFO] Create new WS connection to the server")
	// init commandPromptState new websocket
	commandPromptState.ws = new WebSocketImpl()

	return {
		init: true,
	}
}

