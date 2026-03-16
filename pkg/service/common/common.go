package common

type EchoRequest struct {
	Data string `json: data`
}

type EchoResponse struct {
	Data string `json: data`
}

type EvalRequest struct {
	Data string `json: data`
}

type EvalResponse struct {
	Result string `json: result`
}
