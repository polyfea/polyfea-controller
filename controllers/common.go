package controllers

const (
	PortName = "webserver"
)

func ptr[T any](v T) *T { return &v }
