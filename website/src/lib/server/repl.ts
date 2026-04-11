// JSON schema will be better to handle this
export interface ErrorResponse {
    type: string;
    code: string;
    message: string;
}

export interface EvalRequest {
    data: string;
}

export interface EvalResponseSuccess {
    output?: string;
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
            code: "internal_error",
            type: "internal_error",
            message: 'Internal server error'
        };
    }
}

export function evaluateMock(_req: EvalRequest): EvalResponse {
    return {
        output: 'Mock output'
    };
}
