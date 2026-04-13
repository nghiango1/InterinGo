// JSON schema will be better to handle this
export interface ErrorResponse {
    type: number;
    code: string;
    message: string;
}

// Example
//{
//  "type":400,
//  "code":"parser_error",
//  "message":"Parse error: provided code was invalid,",
//  "error":[
//     {"message":"Expect call expression, but not found closing `)`","range":{"start":{"Line":0,"Character":5},"end":{"Line":0,"Character":7}}}
//   ]
// }
export interface ParserErrorResponse extends ErrorResponse {
    type: 400;
    code: "parser_error";
    error: ParserError[];
}

interface ParserError {
    message: string;
    range: Range;
}

interface Range {
    start: Position;
    end: Position;
}

interface Position {
    start: number;
    end: number;
}

export interface EvalRequest {
    data: string;
}

export interface EvalResponseSuccess {
    output?: string;
}

export type EvalErrorResponse = ErrorResponse | ParserErrorResponse;

export type EvalResponse = EvalResponseSuccess | EvalErrorResponse;

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
            type: 500,
            message: 'Internal server error'
        };
    }
}

export function evaluateMock(_req: EvalRequest): EvalResponse {
    return {
        output: 'Mock output'
    };
}
