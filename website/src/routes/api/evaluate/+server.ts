import type { RequestHandler } from './$types';
import {
	evaluateMock,
	type EvalRequest,
	type EvalResponse
} from '$lib/server/repl';
import { json } from '@sveltejs/kit';

const BAD_REQUEST = 400;
const HTTP_STATUS_OK = 200;

export const POST: RequestHandler = async ({ request }) => {
	let data: EvalRequest;
	try {
		data = await request.json();
	} catch (e) {
		const error: EvalResponse = {
			type: BAD_REQUEST,
			code: BAD_REQUEST.toString(),
			message: 'Input not a valid JSON'
		};

		return json(error, { status: BAD_REQUEST });
	}

	console.log('[INFO] Got eval request: ', data);

	let output: EvalResponse;

	output = evaluateMock(data);
	return json(output, { status: HTTP_STATUS_OK });
};
