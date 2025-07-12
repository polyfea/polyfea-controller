package controllers

import "time"

const (
	timeout  = time.Second * 10
	interval = time.Millisecond * 250
	PortName = "webserver"
)

func ptr[T any](v T) *T { return &v }
