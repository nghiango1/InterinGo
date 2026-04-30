export const prerender = true;
export const ssr = false;

// Keep it here for later
import type { LayoutLoad } from "./$types";

export function load({ }: LayoutLoad) {
	return {
		init: true,
	}
}

