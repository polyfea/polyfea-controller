/*
Copyright 2023.

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
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/controllers"
	"github.com/polyfea/polyfea-controller/repository"
	webserver "github.com/polyfea/polyfea-controller/web-server"
	"github.com/polyfea/polyfea-controller/web-server/configuration"

	//+kubebuilder:scaffold:imports

	"github.com/rs/zerolog"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
	wg       = sync.WaitGroup{}
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(polyfeav1alpha1.AddToScheme(scheme))

	utilruntime.Must(gatewayv1.AddToScheme(scheme))

	utilruntime.Must(networkingv1.AddToScheme(scheme))

	utilruntime.Must(gatewayv1alpha2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a2eec30c.github.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	microFrontendRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.MicroFrontend]()
	microFrontendClassRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.MicroFrontendClass]()
	webComponentRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.WebComponent]()

	if err = (&controllers.MicroFrontendReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Recorder:   mgr.GetEventRecorderFor("microfrontend-controller"),
		Repository: microFrontendRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MicroFrontend")
		os.Exit(1)
	}
	if err = (&controllers.WebComponentReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Recorder:   mgr.GetEventRecorderFor("webcompoent-controller"),
		Repository: webComponentRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "WebComponent")
		os.Exit(1)
	}
	if err = (&controllers.MicroFrontendClassReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Recorder:   mgr.GetEventRecorderFor("microfrontendclass-controller"),
		Repository: microFrontendClassRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MicroFrontendClass")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	initTelemetry(context.Background(), &logger)

	shutdown, err := initTelemetry(ctx, &logger)
	defer shutdown(context.Background())

	if err != nil {
		setupLog.Error(err, "unable to initialize telemetry")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	wg.Add(1)
	go startManager(ctx, cancel, mgr)

	wg.Add(1)
	go startHTTPServer(ctx, cancel, microFrontendClassRepository, microFrontendRepository, webComponentRepository, &logger)

	<-ctx.Done()

	wg.Wait()
}

func startManager(ctx context.Context, cancel context.CancelFunc, mgr manager.Manager) {
	defer wg.Done()
	defer cancel()

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
	}
}

func startHTTPServer(
	ctx context.Context,
	cancel context.CancelFunc,
	microFrontendClassRepository repository.Repository[*polyfeav1alpha1.MicroFrontendClass],
	microFrontendRepository repository.Repository[*polyfeav1alpha1.MicroFrontend],
	webComponentRepository repository.Repository[*polyfeav1alpha1.WebComponent],
	logger *zerolog.Logger) {

	defer wg.Done()
	defer cancel()

	server := &http.Server{
		Addr: ":" + configuration.GetConfigurationValueOrDefault("POLYFEA_WEB_SERVER_PORT", "8082"),
		Handler: otelhttp.NewHandler(
			webserver.SetupRouter(microFrontendClassRepository, microFrontendRepository, webComponentRepository, logger),
			"polyfea-web-server",
		),
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		setupLog.Error(err, "problem running server")
	}
}

func initTelemetry(ctx context.Context, logger *zerolog.Logger) (shutdown func(context.Context) error, err error) {
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
