import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

import type { Component } from 'svelte';

export const load: PageLoad = async ({ params, parent }) => {
    await parent();

    let docs;
    try {
        docs = await import(`$lib/docs/${params.slug}.md`);
    } catch (e) {
        return error(404, `Document not found`);
    }
    let metadata = null;
    if (docs.metadata) {
        metadata = docs.metadata as Record<string, unknown>;
    }
    const content = docs.default as Component;

    return {
        content,
        metadata
    };
};
