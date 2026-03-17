import type { EvalRequest, EvalResponse } from '$lib/server/repl';

export async function postEvaluate(req: EvalRequest): Promise<EvalResponse> {
    const response = await fetch('/api/evaluate', {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            'content-type': 'application/json'
        }
    });

    const output = await response.json();
    return output as EvalResponse;
}
