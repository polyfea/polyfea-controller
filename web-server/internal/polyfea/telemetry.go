package polyfea

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/trace"
)

type instruments struct {
	tracer           trace.Tracer
	meters           metric.Meter
	spa_served       metric.Int64Counter
	boot_served      metric.Int64Counter
	proxied_resource metric.Int64Counter
	context_areas    metric.Int64Counter
	not_found        metric.Int64Counter
}

var telemetry = sync.OnceValue[instruments](func() instruments {
	var err error
	instruments := instruments{}

	instruments.tracer = otel.Tracer("webserver")
	instruments.meters = otel.Meter("webserver")
	instruments.spa_served, err = instruments.meters.Int64Counter(
		"spa_page_served",
		metric.WithDescription("Count of index.html of Single Page App served with success"),
		metric.WithUnit("{count}"),
	)
	if err != nil {
		panic(err)
	}

	instruments.boot_served, err = instruments.meters.Int64Counter(
		"boot_script_served",
		metric.WithDescription("Count of succesfull serving of polyfea boot loader"),
		metric.WithUnit("{count}"),
	)
	if err != nil {
		panic(err)
	}

	instruments.proxied_resource, err = instruments.meters.Int64Counter(
		"proxied_resource",
		metric.WithDescription("Count of requests proxied to backend"),
		metric.WithUnit("{count}"),
	)
	if err != nil {
		panic(err)
	}

	instruments.context_areas, err = instruments.meters.Int64Counter(
		"context_areas",
		metric.WithDescription("Count of requests with context areas"),
		metric.WithUnit("{count}"),
	)
	if err != nil {
		panic(err)
	}

	instruments.not_found, err = instruments.meters.Int64Counter(
		"not_found",
		metric.WithDescription("Count of requests with not found resources"),
		metric.WithUnit("{count}"),
	)

	if err != nil {
		panic(err)
	}
	return instruments
})
