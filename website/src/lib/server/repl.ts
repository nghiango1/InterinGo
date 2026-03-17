// Mock ?
export interface ErrorResponse {
    status: number;
    message: string
}
export interface EvalRequest {
    data: string;
}

export interface EvalResponseSuccess {
    status: 200;
    output: string;
}

export type EvalResponse = EvalResponseSuccess | ErrorResponse;

export function input(req: EvalRequest): EvalResponse {
    return {
        status: 200,
        output: 'hello'
    };
}
