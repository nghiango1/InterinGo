// Keep it here for later
import type { PageLoad } from "./$types";

export function load({} : PageLoad) {
	console.log("[INFO] Loaded page done")
	return {
		init: true,
	}
}

