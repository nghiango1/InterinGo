package core

func (c *ServiceCore) EvaluateHandler(req EvaluateRequest) EvaluateResponse {
	return c.evaluateHandler(ReplRuntime{
		core: c.runtimeCoreDefault,
	}, req.Data)
}
