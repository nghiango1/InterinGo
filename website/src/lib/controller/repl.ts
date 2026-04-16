import type { CreateReplRuntimeRequest, CreateReplRuntimeResponse, EvalRequest, EvalRequestV2, EvalResponse } from '$lib/server/repl';

export async function postEvaluate(req: EvalRequest): Promise<[number, EvalResponse]> {
    const response = await fetch('/api/evaluate', {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            'content-type': 'application/json'
        }
    });

    const output = await response.json();
    return [response.status, output as EvalResponse];
}

export async function postEvaluateV2(req: EvalRequestV2): Promise<[number, EvalResponse]> {
    const response = await fetch(`/api/repl`, {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            'content-type': 'application/json'
        }
    });

    const output = await response.json();
    return [response.status, output as EvalResponse];
}

export async function createReplRuntime(req: CreateReplRuntimeRequest): Promise<[number, CreateReplRuntimeResponse]> {
    const response = await fetch(`/api/repl`, {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            'content-type': 'application/json'
        }
    });

    const output = await response.json();
    return [response.status, output as CreateReplRuntimeResponse];
}
