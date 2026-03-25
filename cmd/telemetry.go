/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func initTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	// prometheus exporter will be in conflict with kubebuilder metrics server
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil, err
	}

	metricProvider :=
		metricsdk.NewMeterProvider(metricsdk.WithReader(metricReader))
	otel.SetMeterProvider(metricProvider)

	traceExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(traceExporter))

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown = func(context.Context) error {
		errMetric := metricProvider.Shutdown(ctx)
		errTrace := traceProvider.Shutdown(ctx)

		if errMetric != nil || errTrace != nil {
			return fmt.Errorf("error shutting down telemetry: %v, %v", errMetric, errTrace)
		}
		return nil
	}

	return shutdown, nil
}
