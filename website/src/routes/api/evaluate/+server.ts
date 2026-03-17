import type { RequestHandler } from './$types';
import type { EvalRequest, EvalResponse } from '$lib/server/repl';
import { json } from '@sveltejs/kit';

const BAD_REQUEST = 400
const OK = 200

export const POST: RequestHandler = async ({ request }) => {
	let data: EvalRequest ;
	try {
		const { done } = await request.json();
		data = done
	} catch (e) {
		const error: EvalResponse = {
			status: BAD_REQUEST,
			message: 'Input not a valid JSON'
		};

		return json(error, { status: BAD_REQUEST })
	}

	console.log(data)

	const output: EvalResponse = {
		status: OK,
		output: 'Not implemented'
	};

	return json(output, { status: 200 });
};
