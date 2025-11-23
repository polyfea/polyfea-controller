package controller

const (
	PortName                 = "webserver"
	DefaultFrontendClassName = "polyfea-controller-default"
)

func ptr[T any](v T) *T { return &v }
