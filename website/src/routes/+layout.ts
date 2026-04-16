export const prerender = true;
export const ssr = false;

// This help with building index file, as we have nested directory
// that can cause different redirect behavior to static file hosting
// - Directory `docs/` is forced redirected to 'docs/index.html'
// - No directory 'docs.html' is used as literal file, golang doesn't support
// route it to /docs
export const trailingSlash = 'always';

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

