export interface ErrorResponse {
    status: 400 | 500;
    message: string;
}
export interface EvalRequest {
    data: string;
}

export interface EvalResponseSuccess {
    status: 200;
    output: string;
}

export type EvalResponse = EvalResponseSuccess | ErrorResponse;

export async function evaluateServer(req: EvalRequest, base?: string): Promise<EvalResponse> {
    const response = await fetch(new URL('/api/evaluate', base), {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            'content-type': 'application/json'
        }
    });

    try {
        const output = await response.json();
        return output as EvalResponse;
    } catch (e) {
        return {
            status: 500,
            message: 'Server error'
        };
    }
}

export function evaluateMock(_req: EvalRequest): EvalResponse {
    return {
        status: 200,
        output: 'Mock output'
    };
}
