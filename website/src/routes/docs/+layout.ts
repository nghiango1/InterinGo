import type { DocInfo, NavigationRecord } from '$lib/type';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	let navSessions: NavigationRecord[] = [];
	const keys: Map<string, number> = new Map<string, number>();
	const raw = import.meta.env.VITE_DOCS;
	if (!raw) {
		console.warn('[navSessions] VITE_DOCS_MANIFEST not set – did you run the build?');
	}
	try {
		const docInfos = JSON.parse(raw) as DocInfo[];
		for (let i = 0; i < docInfos.length; i++) {
			const session = docInfos[i].session || 'Misc';
			if (!keys.has(session)) {
				keys.set(session, navSessions.length);
				navSessions.push({
					name: session,
					docs: []
				});
			}
			navSessions[keys.get(session)!].docs.push(docInfos[i]);
		}
	} catch {
		console.error('[navSessions] Failed to parse VITE_DOCS_MANIFEST');
	}

	return {
		navSessions: navSessions,
	};
};
