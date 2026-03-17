import type { RequestHandler } from './$types';
import {
	evaluateMock,
	evaluateServer,
	type EvalRequest,
	type EvalResponse
} from '$lib/server/repl';
import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const BAD_REQUEST = 400;
const HTTP_STATUS_OK = 200;

export const POST: RequestHandler = async ({ request }) => {
	let data: EvalRequest;
	try {
		const { done } = await request.json();
		data = done;
	} catch (e) {
		const error: EvalResponse = {
			status: BAD_REQUEST,
			message: 'Input not a valid JSON'
		};

		return json(error, { status: BAD_REQUEST });
	}

	console.log('[INFO] Input: ', data);

	let output: EvalResponse;

	if (env.BACKEND_SERVER_URL) {
		output = await evaluateServer(data, env.BACKEND_SERVER_URL);
	} else {
		output = evaluateMock(data);
	}

	return json(output, { status: HTTP_STATUS_OK });
};
